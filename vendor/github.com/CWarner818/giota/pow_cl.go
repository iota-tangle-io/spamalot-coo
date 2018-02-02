// +build gpu

/*
MIT License

Copyright (c) 2017 Shinya Yagyu

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package giota

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync/atomic"
	"unsafe"

	"github.com/jgillich/go-opencl/cl"
)

var loopcount byte = 32

var countCL int64

type bufferInfo struct {
	size    int64
	flag    cl.MemFlag
	isArray bool
	data    []byte
}

func init() {
	// TODO: update to Curl-P-81
	// pows["PowCL"] = PowCL
}

var stopCL = true

// nolint: gocyclo
func exec(
	que *cl.CommandQueue,
	ker []*cl.Kernel,
	cores, nlocal int,
	mobj []*cl.MemObject,
	founded *int32,
	tryte chan Trytes) error {

	// initialize
	nglobal := cores * nlocal
	ev1, err := que.EnqueueNDRangeKernel(ker[0], nil, []int{nglobal}, []int{nlocal}, nil)
	if err != nil {
		return err
	}
	defer ev1.Release()

	if err = que.Finish(); err != nil {
		return err
	}

	found := make([]byte, 1)
	var cnt int
	num := int64(cores) * 64 * int64(loopcount)
	for cnt = 0; found[0] == 0 && *founded == 0 && !stopCL; cnt++ {
		// start searching
		var ev2 *cl.Event
		ev2, err = que.EnqueueNDRangeKernel(ker[1], nil, []int{nglobal}, []int{nlocal}, nil)
		if err != nil {
			return err
		}
		defer ev2.Release()
		dataPtr := unsafe.Pointer(&found[0])
		dataSize := int(unsafe.Sizeof(found[0])) * len(found)
		ev3, err := que.EnqueueReadBuffer(mobj[6], true, 0, dataSize, dataPtr, []*cl.Event{ev2})
		if err != nil {
			return err
		}
		ev3.Release()

		atomic.AddInt64(&countCL, num)
	}

	if *founded != 0 || stopCL {
		return nil
	}

	atomic.StoreInt32(founded, 1)

	// finalize, get the result.
	ev4, err := que.EnqueueNDRangeKernel(ker[2], nil, []int{nglobal}, []int{nlocal}, nil)
	if err != nil {
		return err
	}
	defer ev4.Release()

	result := make([]byte, HashSize*8)
	dataPtr := unsafe.Pointer(&result[0])
	dataSize := int(unsafe.Sizeof(result[0])) * len(result)
	ev5, err := que.EnqueueReadBuffer(mobj[0], true, 0, dataSize, dataPtr, []*cl.Event{ev4})
	if err != nil {
		return err
	}
	ev5.Release()

	rr := make(Trits, HashSize)
	for i := 0; i < HashSize; i++ {
		switch {
		case result[i*8] == 0xff:
			rr[i] = -1
		case result[i*8] == 0x0 && result[i*8+7] == 0x0:
			rr[i] = 0
		case result[i*8] == 0x1 || result[i*8+7] == 0x1:
			rr[i] = 1
		}
	}

	tryte <- rr.Trytes()
	return nil
}

// nolint: gocyclo
func loopCL(binfo []bufferInfo) (Trytes, error) {
	defers := make([]func(), 0, 10)
	defer func() {
		for _, f := range defers {
			f()
		}
	}()

	platforms, err := cl.GetPlatforms()
	if err != nil {
		return "", err
	}

	exist := false
	var founded int32
	result := make(chan Trytes)
	for _, p := range platforms {
		var devs []*cl.Device
		devs, err = p.GetDevices(cl.DeviceTypeGPU)
		if err != nil || len(devs) == 0 {
			continue
		}

		exist = true
		// TODO: this case checks the error after appending, but all the other cases below
		// do it the opposite way. Check to see if this can be reversed to maintain the
		// pattern
		cont, err := cl.CreateContext(devs)
		defers = append(defers, cont.Release)
		if err != nil {
			return "", err
		}

		prog, err := cont.CreateProgramWithSource([]string{kernel})
		if err != nil {
			return "", err
		}

		defers = append(defers, prog.Release)
		if err := prog.BuildProgram(devs, "-Werror"); err != nil {
			println(p.Name())
			return "", err
		}

		ker := make([]*cl.Kernel, 3)
		defers = append(defers, func() {
			for _, k := range ker {
				if k != nil {
					k.Release()
				}
			}
		})

		for i, n := range []string{"init", "search", "finalize"} {
			ker[i], err = prog.CreateKernel(n)
			if err != nil {
				return "", err
			}
		}

		for _, d := range devs {
			mult := d.MaxWorkGroupSize()
			cores := d.MaxComputeUnits()
			mmax := d.MaxMemAllocSize()
			isLittle := d.EndianLittle()
			nlocal := 0
			for nlocal = stateSize; nlocal > mult; {
				nlocal /= 3
			}

			var totalmem int64
			mobj := make([]*cl.MemObject, len(binfo))
			defers = append(defers, func() {
				for _, o := range mobj {
					if o != nil {
						o.Release()
					}
				}
			})

			que, err := cont.CreateCommandQueue(d, 0)
			if err != nil {
				return "", err
			}

			defers = append(defers, que.Release)

			for i, inf := range binfo {
				msize := inf.size
				if inf.isArray {
					msize *= int64(cores * mult)
				}

				if totalmem += msize; totalmem > mmax {
					//return "", errors.New("max memory passed")
					fmt.Println("max memory passed")
				}

				mobj[i], err = cont.CreateEmptyBuffer(inf.flag, int(msize))
				if err != nil {
					return "", err
				}

				if inf.data != nil {
					var ev *cl.Event
					switch {
					case isLittle:
						dataPtr := unsafe.Pointer(&inf.data[0])
						dataSize := int(unsafe.Sizeof(inf.data[0])) * len(inf.data)
						ev, err = que.EnqueueWriteBuffer(mobj[i], true, 0, dataSize, dataPtr, nil)
					default:
						data := make([]byte, len(inf.data))
						for i := range inf.data {
							data[i] = inf.data[len(inf.data)-i-1]
						}
						dataPtr := unsafe.Pointer(&data[0])
						dataSize := int(unsafe.Sizeof(data[0])) * len(data)
						ev, err = que.EnqueueWriteBuffer(mobj[i], true, 0, dataSize, dataPtr, nil)
					}

					if err != nil {
						return "", err
					}
					ev.Release()
				}

				for _, k := range ker {
					if err := k.SetArg(i, mobj[i]); err != nil {
						return "", err
					}
				}
			}

			go func() {
				err := exec(que, ker, cores, nlocal, mobj, &founded, result)
				if err != nil {
					panic(err)
				}
			}()
		}
	}

	if !exist {
		return "", errors.New("no GPU found")
	}

	r := <-result
	close(result)

	return r, nil
}

// PowCL is proof of work of iota in OpenCL.
func PowCL(trytes Trytes, mwm int) (Trytes, error) {
	switch {
	case !stopCL:
		stopCL = true
		return "", errors.New("pow is already running, stopped")
	case trytes == "":
		return "", errors.New("invalid trytes")
	}

	stopCL = false
	countCL = 0
	c := NewCurl()
	c.Absorb(trytes[:(transactionTrinarySize-HashSize)/3])
	tr := trytes.Trits()
	copy(c.state, tr[transactionTrinarySize-HashSize:])

	lmid, hmid := para(c.state)
	lmid[0] = low0
	hmid[0] = high0
	lmid[1] = low1
	hmid[1] = high1
	lmid[2] = low2
	hmid[2] = high2
	lmid[3] = low3
	hmid[3] = high3

	low := make([]byte, 8*stateSize)
	for i, v := range lmid {
		binary.LittleEndian.PutUint64(low[8*i:], v)
	}

	high := make([]byte, 8*stateSize)
	for i, v := range hmid {
		binary.LittleEndian.PutUint64(high[8*i:], v)
	}

	binfo := []bufferInfo{
		bufferInfo{
			8 * HashSize, cl.MemWriteOnly, false, nil,
		},
		bufferInfo{
			8 * stateSize, cl.MemReadWrite, true, low, //mid_low
		},
		bufferInfo{
			8 * stateSize, cl.MemReadWrite, true, high, //mid_high
		},
		bufferInfo{
			8 * stateSize, cl.MemReadWrite, true, nil,
		},
		bufferInfo{
			8 * stateSize, cl.MemReadWrite, true, nil,
		},
		bufferInfo{
			8, cl.MemWriteOnly, false, []byte{byte(mwm), 0, 0, 0, 0, 0, 0, 0}, // mwm
		},
		bufferInfo{
			1, cl.MemReadWrite, false, nil,
		},
		bufferInfo{
			8, cl.MemReadWrite, true, nil,
		},
		bufferInfo{
			8, cl.MemWriteOnly, false, []byte{loopcount, 0, 0, 0, 0, 0, 0, 0}, //loop_count
		},
	}

	return loopCL(binfo)
}
