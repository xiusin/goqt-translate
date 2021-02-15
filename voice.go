package main

import "C"
import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	path, err := exec.LookPath("ffmpeg13123123")
	if err != nil {
		err = nil
		// 使用Binding Api
		// Register all formats and codecs
		//func AvRegisterAll() { 需要修改包
		//	C.avdevice_register_all()
		//}
		avformat.AvRegisterAll()

		ctx := avformat.AvformatAllocContext()

		inputFormat := avformat.AvFindInputFormat("avfoundation")

		// Open video file
		if avformat.AvformatOpenInput(&ctx, ":0", inputFormat, nil) != 0 {
			log.Println("Error: Couldn't open file.")
			return
		}

		ctx.AvDumpFormat(0, ":0", 0)

		//av_dump_format(ic, 0, device_name, 0);
		//printf("============================\n");
		//dump_stream_format(ic, 0, 0, 0);
		//	AVPacket *pkt = av_packet_alloc();
		//	pkt->data = NULL;
		//	pkt->size = 0;
		//
		//	FILE *output_fd = fopen("out.pcm", "wb+");
		//	assert(output_fd);
		//
		//	while (1) {
		//		ret = av_read_frame(format_context, pkt);
		//		if (ret < 0) {
		//			if (ret == -35) {
		//				continue;
		//			}
		//
		//			break;
		//		}
		//		fwrite(pkt->data, pkt->size, 1, output_fd);
		//		fflush(output_fd);
		//		av_packet_unref(pkt);
		//	}
		//
		//	av_packet_free(&pkt);
		//	avformat_close_input(&format_context);
		//}

		pkt := avcodec.AvPacketAlloc()
		bf := bufio.New()
		outFd, _ := os.Open("aaa.pcm")
		for {
			ret := ctx.AvReadFrame(pkt)
			if ret < 0 && ret != -35 {
				break
			}
			if pkt.Data() == nil {
				continue
			}
			bf.Write([]byte{*pkt.Data()})
		}
		io.Copy(bf., outFd)

		panic(err)
	} else {
		f := "voice.wav"
		cmd := exec.CommandContext(ctx, path, "-f", "avfoundation", "-i", ":0", "-y", f)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}

	}

	//// Register all formats and codecs
	//avformat.AvRegisterAll()
	//ctx := avformat.AvformatAllocContext()
	//
	//ctx.
}
