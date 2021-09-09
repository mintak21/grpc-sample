package service

import (
	"context"
	"math/rand"
	"sync"
	"time"

	ts "github.com/golang/protobuf/ptypes/timestamp"

	pb "github.com/mintak21/proto/sample/golang"
)

const (
	chef_name = "Tom Brown"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generate service
func NewBakePancakeService() *BakePancakeService {
	return &BakePancakeService{
		report: &report{
			data: make(map[pb.Pancake_Menu]int),
		},
	}
}

type BakePancakeService struct {
	report *report
}

type report struct {
	mx   sync.Mutex
	data map[pb.Pancake_Menu]int
}

func (s *BakePancakeService) Bake(ctx context.Context, request *pb.BakeRequest) (response *pb.BakeResponse, err error) {
	s.report.mx.Lock()
	s.report.data[request.Menu] = s.report.data[request.Menu] + 1
	s.report.mx.Unlock()
	now := time.Now()

	return &pb.BakeResponse{
		Pancake: &pb.Pancake{
			BakerName: chef_name,
			Menu:      request.Menu,
			CreateTime: &ts.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

func (s *BakePancakeService) Report(ctx context.Context, request *pb.ReportRequest) (response *pb.ReportResponse, err error) {
	counts := make([]*pb.Report_BakeCount, 0)

	s.report.mx.Lock()
	for menu, count := range s.report.data {
		counts = append(counts, &pb.Report_BakeCount{
			Menu:  menu,
			Count: int32(count),
		})
	}
	s.report.mx.Unlock()

	return &pb.ReportResponse{
		Report: &pb.Report{
			BakeCounts: counts,
		},
	}, nil
}
