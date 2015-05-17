package dev

import (
	"code.google.com/p/portaudio-go/portaudio"
	"github.com/edmccard/avr-sim/core"
)

type Speaker struct {
	stream       *portaudio.Stream
	channel      chan float32
	timer        *core.Timer
	curSample    float32
	avgBuf       []float32
	avgIdx       int
	cycPerSample uint
	pin          byte
	mask         byte
	lastToggle   int64
	started      bool
}

func NewSpeaker(timer *core.Timer, hertz uint, pin int) (*Speaker, error) {
	spk := &Speaker{curSample: -1.0}
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	parameters.Output.Channels = 1
	parameters.SampleRate = 44100
	stream, err := portaudio.OpenStream(parameters, spk.Callback)
	if err != nil {
		return nil, err
	}
	spk.timer = timer
	spk.cycPerSample = hertz / uint(parameters.SampleRate)
	spk.channel = make(chan float32, 8192)
	spk.avgBuf = make([]float32, spk.cycPerSample)
	spk.mask = 1 << uint(pin)
	spk.stream = stream
	return spk, nil
}

func (spk *Speaker) OnSlice() error {
	if !spk.started {
		spk.started = true
		return spk.Start()
	}
	return nil
}

func (spk *Speaker) Start() error {
	return spk.stream.Start()
}

func (spk *Speaker) Stop() error {
	return spk.stream.Stop()
}

func (spk *Speaker) Callback(out []float32) {
	for i := range out {
		select {
		case sample := <-spk.channel:
			out[i] = sample
		default:
			out[i] = 0
		}
	}
}

func (spk *Speaker) Write(addr core.Addr, val byte) {
	val &= spk.mask
	if val == spk.pin {
		return
	}
	spk.pin = val
	spk.makeSamples()
	spk.curSample *= -1.0
}

func (spk *Speaker) makeSamples() {
	curCycle := spk.timer.GetCount()
	elapsed := curCycle - spk.lastToggle
	spk.lastToggle = curCycle

	if spk.avgIdx != 0 {
		for ; spk.avgIdx < len(spk.avgBuf); spk.avgIdx++ {
			if elapsed == 0 {
				break
			}
			spk.avgBuf[spk.avgIdx] = spk.curSample
			elapsed--
		}
		if spk.avgIdx == len(spk.avgBuf) {
			spk.avgIdx = 0
			avg := float32(0.0)
			for _, sample := range spk.avgBuf {
				avg += sample
			}
			avg /= float32(len(spk.avgBuf))
			spk.sendSample(avg)
		}
	}

	for i := int64(0); i < elapsed/int64(spk.cycPerSample); i++ {
		spk.sendSample(spk.curSample)
	}

	avgCycs := elapsed % int64(spk.cycPerSample)
	if avgCycs != 0 {
		for spk.avgIdx = 0; spk.avgIdx < int(avgCycs); spk.avgIdx++ {
			spk.avgBuf[spk.avgIdx] = spk.curSample
		}
	}
}

func (spk *Speaker) sendSample(sample float32) {
	select {
	case spk.channel <- sample:
	default:
	}
}
