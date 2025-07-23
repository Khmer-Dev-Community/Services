package videosnap

import (
	"encoding/json"
	"fmt"
	"net/http"
	redis "telegram-service/config"
	"telegram-service/lib/desks"
	"telegram-service/lib/videosnap/setting"
	"telegram-service/utils"

	"time"

	"github.com/sirupsen/logrus"
)

// VideosnapServiceInterface defines the methods for the VideosnapService
type VideosnapServiceInterface interface {
	VideosnapServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse
	VideosnapServiceCreate(obj []byte) (*VideosnapResponseDTO, error)
	VideosnapServiceUpdate(obj []byte) (*VideosnapResponseDTO, error)
	VideosnapServiceGetById(obj []byte) (*VideosnapResponseDTO, error)

	VideosnapServiceDelete(obj []byte) (bool, error)
}
type VideosnapService struct {
	repo *VideosnapRepository
}

func NewVideosnapService(repo *VideosnapRepository) *VideosnapService {
	return &VideosnapService{repo: repo}
}

// VideosnapServiceGetList retrieves a paginated list of videosnaps
func (vs *VideosnapService) VideosnapServiceGetList(requestDto utils.PaginationRequestDTO) utils.ServicePaginationResponse {
	if requestDto.Page <= 0 {
		requestDto.Page = 1 // Default to first page if invalid
	}
	if requestDto.Limit <= 0 {
		requestDto.Limit = 10 // Default limit if invalid
	}
	videosnaps, total, err := vs.repo.GetVideosnapList(int(requestDto.Page), int(requestDto.Limit), requestDto.Query)
	if err != nil {
		utils.ErrorLog(nil, err.Error())                                                                                                                               // Log the error
		return utils.NewServicePaginationResponse(nil, 0, 0, 0, http.StatusInternalServerError, err.Error(), logrus.ErrorLevel, "VideosnapService [GetVideosnapList]") // Return error response
	}
	videosnapDTOs := make([]*VideosnapResponseDTO, len(videosnaps))
	for i, video := range videosnaps {
		videoDTO := &VideosnapResponseDTO{}
		videoDTO.FromModel(video)
		videosnapDTOs[i] = videoDTO
	}
	return utils.NewServicePaginationResponse(videosnapDTOs, total, int(requestDto.Page), int(requestDto.Limit), http.StatusOK, string(utils.SuccessMessage), logrus.InfoLevel, "VideosnapService [GetVideosnapList]")
}

// VideosnapServiceCreate creates a new videosnap
func (vs *VideosnapService) VideosnapServiceCreate(obj []byte) (*VideosnapResponseDTO, error) {
	var createDto VideosnapCreateDTO
	err := json.Unmarshal(obj, &createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	// Lookup in Table Setting if set
	redisDeskKey := fmt.Sprintf("deskSetting:table:%s", createDto.Streamkey)
	tableSetting, err := redis.Get(redisDeskKey)
	if err != nil {
		utils.ErrorLog(fmt.Errorf("Failed to get data from Redis: %s", redisDeskKey), err.Error())
	}
	var tableSnapSetting desks.DeskSettingResponseDTO

	if err := json.Unmarshal([]byte(tableSetting), &tableSnapSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %s", redisDeskKey), err.Error())
	}
	utils.InfoLog(tableSnapSetting, string(utils.SuccessMessage))
	var Server_ID uint // Change to appropriate type if necessary
	if tableSnapSetting.DeskStreamServer != 0 {
		Server_ID = tableSnapSetting.DeskStreamServer
	} else {
		Server_ID = 1 // Default Server ID if StreamServer.ID is not set
	}
	utils.InfoLog(Server_ID, string(utils.SuccessMessage))

	redisServerKey := fmt.Sprintf("videoSetting:server:%d", Server_ID)
	videoInfo, err := redis.Get(redisServerKey)
	if err != nil {
		utils.ErrorLog(fmt.Errorf("failed to get video settings from Redis: %v", err), err.Error())
	}
	utils.InfoLog(fmt.Sprintf("videoSetting:server:%d", Server_ID), string(utils.SuccessMessage))
	var videoSnapSetting setting.VideoSnapSettingResponseDTO
	if err := json.Unmarshal([]byte(videoInfo), &videoSnapSetting); err != nil {
		utils.ErrorLog(fmt.Errorf("failed to unmarshal value from Redis: %v", err), err.Error())
	}
	utils.InfoLog(videoSnapSetting, string(utils.SuccessMessage))
	processId, err := SnapVideo(createDto.GameNo, createDto.Streamkey)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to start video snap: %v", err)
	}

	createDto.ProcessId = uint(processId)
	createDto.Rtmp = videoSnapSetting.Rtmp
	createDto.StorePath = videoSnapSetting.OutputPath
	currentDate := time.Now().Format("20060102")

	VideoLink := fmt.Sprintf("%s/%s/%s/%s.mp4",
		videoSnapSetting.AcceessDomain,
		currentDate,
		createDto.Streamkey,
		createDto.Streamkey+createDto.GameNo,
	)
	ImageLink := fmt.Sprintf("%s/%s/%s/%s.jpeg",
		videoSnapSetting.AcceessDomain,
		currentDate,
		createDto.Streamkey,
		createDto.Streamkey+createDto.GameNo,
	)
	createDto.ImageURL = ImageLink
	createDto.VideoURL = VideoLink
	// Create the videosnap in the repository
	createdVideosnap, err := vs.repo.CreateVideosnap(createDto)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	// Task save redis
	redisKey := fmt.Sprintf("task:snapshoot:%s", createdVideosnap.GameNo)
	expiration := time.Minute * 5 // 1 nimutes
	// Store in Redis without expiration
	if err := redis.SetWithExpiration(redisKey, createdVideosnap, expiration); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("failed to store setting in Redis: %v", err)
	}
	// Create the response DTO
	videosnapDTO := &VideosnapResponseDTO{}
	videosnapDTO.FromModel(createdVideosnap)

	utils.InfoLog(videosnapDTO, string(utils.SuccessMessage))
	return videosnapDTO, nil
}

// VideosnapServiceUpdate updates an existing videosnap
func (vs *VideosnapService) VideosnapServiceUpdate(obj []byte) (*VideosnapResponseDTO, error) {
	var updateVideosnapRequest VideosnapUpdateDTO
	if err := json.Unmarshal(obj, &updateVideosnapRequest); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	videosnap, err := vs.repo.GetVideosnapByGameNo(*updateVideosnapRequest.GameNo)
	if err != nil {
		utils.LoggerRepository(err, "GetVideosnapByGameNo")
		return nil, fmt.Errorf("failed to retrieve videosnap: %v", err)
	}
	// append DATA
	status := true
	updateVideosnapRequest.Status = &status
	SnapShootImage(videosnap.GameNo, videosnap.Streamkey)
	updatedVideosnap, err := vs.repo.UpdateVideosnap(videosnap.ID, updateVideosnapRequest)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}

	updatedDto := &VideosnapResponseDTO{}
	updatedDto.FromModel(updatedVideosnap)
	utils.InfoLog(updatedDto, string(utils.SuccessMessage))
	return updatedDto, nil
}

// VideosnapServiceGetById retrieves a videosnap by ID
func (vs *VideosnapService) VideosnapServiceGetById(obj []byte) (*VideosnapResponseDTO, error) {
	var videosnapDTO VideosnapUpdateDTO
	if err := json.Unmarshal(obj, &videosnapDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	videosnap, err := vs.repo.GetVideosnapByID(videosnapDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return nil, fmt.Errorf("%v", err)
	}
	resultDto := &VideosnapResponseDTO{}
	resultDto.FromModel(videosnap)
	utils.InfoLog(resultDto, string(utils.SuccessMessage))
	return resultDto, nil
}

// VideosnapServiceDelete deletes a videosnap by ID
func (vs *VideosnapService) VideosnapServiceDelete(obj []byte) (bool, error) {
	var videosnapDTO VideosnapResponseDTO
	if err := json.Unmarshal(obj, &videosnapDTO); err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}
	success, err := vs.repo.DeleteVideosnapByID(videosnapDTO.ID)
	if err != nil {
		utils.ErrorLog(nil, err.Error())
		return false, fmt.Errorf("%v", err)
	}

	utils.InfoLog(videosnapDTO, string(utils.SuccessMessage))
	return success, nil
}
