using AutoMapper;
using FluentAssertions;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using NSubstitute;
using UserService;
using Webby.VideoService.ComplaintGrpcClient;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Helpers;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models.Enums;

namespace Webby.FunctionalTests.Tests.Video;

public class VideoServiceTests
{
    private readonly IVideoRepository _videoRepositoryMock;
    private readonly IStorageService _storageServiceMock;
    private readonly ComplaintGrpcService.ComplaintGrpcServiceClient _complaintGrpcClientMock;
    private readonly ITagService _tagServiceMock;
    private readonly GrpcClients.UserService.UserGrpcService.UserGrpcServiceClient _userClientMock;
    private readonly IMapper _mapperMock;
    private readonly IBackgroundTaskQueue _queueMock;
    private readonly IServiceScopeFactory _scopeFactoryMock;
    private readonly IYouTubeSearchService _youtubeSearchServiceMock;
    private readonly ITwitchSearchService _twitchSearchServiceMock;
    private readonly IEventPublisher _eventPublisherMock;
    private readonly ILogger<VideoService.Services.VideoService> _loggerMock;

    private readonly VideoService.Services.VideoService _sut;

    public VideoServiceTests()
    {
        _videoRepositoryMock = Substitute.For<IVideoRepository>();
        _storageServiceMock = Substitute.For<IStorageService>();
        _complaintGrpcClientMock = Substitute.For<ComplaintGrpcService.ComplaintGrpcServiceClient>();
        _tagServiceMock = Substitute.For<ITagService>();
        _userClientMock = Substitute.For<GrpcClients.UserService.UserGrpcService.UserGrpcServiceClient>();
        _mapperMock = Substitute.For<IMapper>();
        _queueMock = Substitute.For<IBackgroundTaskQueue>();
        _scopeFactoryMock = Substitute.For<IServiceScopeFactory>();
        _youtubeSearchServiceMock = Substitute.For<IYouTubeSearchService>();
        _twitchSearchServiceMock = Substitute.For<ITwitchSearchService>();
        _eventPublisherMock = Substitute.For<IEventPublisher>();
        _loggerMock = Substitute.For<ILogger<VideoService.Services.VideoService>>();

        _sut = new VideoService.Services.VideoService(
            _videoRepositoryMock,
            _storageServiceMock,
            _tagServiceMock,
            _userClientMock,
            _mapperMock,
            _queueMock,
            _scopeFactoryMock,
            _youtubeSearchServiceMock,
            _twitchSearchServiceMock,
            _eventPublisherMock,
            _loggerMock,
            _complaintGrpcClientMock
        );
    }

    #region GetVideoById Tests

    [Fact]
    public async Task GetVideoById_WhenVideoNotFound_ShouldThrowApiException404()
    {
        // Arrange
        var videoId = Guid.NewGuid();
        _videoRepositoryMock.FindById(videoId).Returns((VideoService.Models.Video)null);

        // Act
        Func<Task> act = async () => await _sut.GetVideoById(videoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 404 && e.Message == "Get video error");
    }

    [Fact]
    public async Task GetVideoById_WhenVideoIsPrivateAndShowPrivateIsFalse_ShouldThrowApiException403()
    {
        // Arrange
        var videoId = Guid.NewGuid();
        var video = new VideoService.Models.Video { Name = "Some video", VideoId = videoId, IsPrivate = true };
        _videoRepositoryMock.FindById(videoId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.GetVideoById(videoId, showPrivate: false);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Get video error");
    }

    [Fact]
    public async Task GetVideoById_WhenVideoIsBanned_ShouldThrowApiException400()
    {
        // Arrange
        var videoId = Guid.NewGuid();
        var video = new VideoService.Models.Video { Name = "Some video", VideoId = videoId, IsPrivate = false, IsBanned = true };
        _videoRepositoryMock.FindById(videoId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.GetVideoById(videoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 400 && e.Message == "Get video error");
    }

    [Fact]
    public async Task GetVideoById_WhenVideoNotUploaded_ShouldThrowApiException400()
    {
        // Arrange
        var videoId = Guid.NewGuid();
        var video = new VideoService.Models.Video 
        {
            Name = "Some video",
            VideoId = videoId, 
            IsPrivate = false, 
            IsBanned = false, 
            VideoUploadStatus = VideoStatus.Uploading 
        };
        _videoRepositoryMock.FindById(videoId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.GetVideoById(videoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 400 && e.Message == "Get video error");
    }

    [Fact]
    public async Task GetVideoById_WhenValidAndPublished_ShouldReturnVideo()
    {
        // Arrange
        var videoId = Guid.NewGuid();
        var video = new VideoService.Models.Video 
        { 
            Name = "Some video",
            VideoId = videoId, 
            IsPrivate = false, 
            IsBanned = false, 
            VideoUploadStatus = VideoStatus.Ready,
            IsPublished = true
        };
        _videoRepositoryMock.FindById(videoId).Returns(video);

        // Act
        var result = await _sut.GetVideoById(videoId);

        // Assert
        result.Should().NotBeNull();
        result.VideoId.Should().Be(videoId);
    }

    #endregion

    #region DeleteVideo Tests

    [Fact]
    public async Task DeleteVideo_WhenVideoNotFound_ShouldThrowApiException404()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";
        
        _videoRepositoryMock.FindById(actualId).Returns((VideoService.Models.Video)null);

        // Act
        Func<Task> act = async () => await _sut.DeleteVideo(userId, prefixedVideoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 404 && e.Message == "Delete video error");
    }

    [Fact]
    public async Task DeleteVideo_WhenUserIsNotOwner_ShouldThrowApiException403()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var ownerId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video { Name = "Some video", VideoId = actualId, UserId = ownerId };
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.DeleteVideo(requestUserId, prefixedVideoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Delete video error");

        await _videoRepositoryMock.DidNotReceive().DeleteAsync(Arg.Any<Guid>());
    }

    [Fact]
    public async Task DeleteVideo_WhenValid_ShouldDeleteFilesAndRepositoryRecord()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video 
        { 
            Name = "Some video",
            VideoId = actualId, 
            UserId = userId,
            VideoUrl = "some_video_url",
            PreviewUrl = "some_preview_url"
        };

        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        await _sut.DeleteVideo(userId, prefixedVideoId);

        // Assert
        await _storageServiceMock.Received(1).DeleteFileAsync("some_video_url");
        await _storageServiceMock.Received(1).DeleteFileAsync("some_preview_url");
        await _videoRepositoryMock.Received(1).DeleteAsync(actualId);
    }

    #endregion
    
    #region CancelVideoUploading Tests

    [Fact]
    public async Task CancelVideoUploading_WhenUserIsNotOwner_ShouldThrowApiException403()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var ownerId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video { Name = "Some video", VideoId = actualId, UserId = ownerId };
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.CancelVideoUploading(requestUserId, prefixedVideoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Publish video error");
    }

    [Fact]
    public async Task CancelVideoUploading_WhenStatusIsUploading_ShouldSetStatusToCanceled()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video 
        { 
            Name = "Some video",
            VideoId = actualId, 
            UserId = userId, 
            VideoUploadStatus = VideoStatus.Uploading 
        };
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        await _sut.CancelVideoUploading(userId, prefixedVideoId);

        // Assert
        video.VideoUploadStatus.Should().Be(VideoStatus.Canceled);
        await _videoRepositoryMock.Received(1).Update(video);
        await _videoRepositoryMock.DidNotReceive().DeleteAsync(Arg.Any<Guid>());
    }

    [Fact]
    public async Task CancelVideoUploading_WhenStatusIsReady_ShouldDeleteFileAndRecord()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video 
        {
            Name = "Some video",
            VideoId = actualId, 
            UserId = userId, 
            VideoUploadStatus = VideoStatus.Ready,
            VideoUrl = "some_url"
        };
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        await _sut.CancelVideoUploading(userId, prefixedVideoId);

        // Assert
        await _storageServiceMock.Received(1).DeleteFileAsync("some_url");
        await _videoRepositoryMock.Received(1).DeleteAsync(actualId);
    }

    #endregion

    #region CheckUploadStatus Tests

    [Fact]
    public async Task CheckUploadStatus_WhenVideoNotFound_ShouldReturnFalse()
    {
        // Arrange
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";
        _videoRepositoryMock.FindById(actualId).Returns((VideoService.Models.Video)null);

        // Act
        var result = await _sut.CheckUploadStatus(prefixedVideoId);

        // Assert
        result.Should().BeFalse();
    }

    [Theory]
    [InlineData(VideoStatus.Uploading, false)]
    [InlineData(VideoStatus.Canceled, false)]
    [InlineData(VideoStatus.Failed, false)]
    [InlineData(VideoStatus.Ready, true)]
    public async Task CheckUploadStatus_ShouldReturnCorrectStatus(VideoStatus currentStatus, bool expectedResult)
    {
        // Arrange
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";
        var video = new VideoService.Models.Video {Name = "Some video",VideoId = actualId, VideoUploadStatus = currentStatus };
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        var result = await _sut.CheckUploadStatus(prefixedVideoId);

        // Assert
        result.Should().Be(expectedResult);
    }

    #endregion

    #region UpdateVideoInformation Tests

    [Fact]
    public async Task UpdateVideoInformation_WhenUserIsNotOwner_ShouldThrowApiException403()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var ownerId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        
        var request = new Webby.VideoService.Dtos.Video.UpdateVideoRequest 
        { 
            VideoId = $"wb_{actualId}", 
            Name = "New Name" 
        };

        var video = new VideoService.Models.Video { Name = "Some video",VideoId = actualId, UserId = ownerId };
        _videoRepositoryMock.GetVideoInformationById(actualId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.UpdateVideoInformation(requestUserId, request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 403 && e.Message == "Update information error");
    }

    [Fact]
    public async Task UpdateVideoInformation_WhenValidWithoutPreview_ShouldUpdateProperties()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        
        var request = new Webby.VideoService.Dtos.Video.UpdateVideoRequest 
        { 
            VideoId = $"wb_{actualId}", 
            Name = "Updated Name",
            Description = "Updated Desc",
            IsPrivate = true
        };

        var video = new VideoService.Models.Video 
        {
            Name = "Some video",
            VideoId = actualId, 
            UserId = userId,
            IsPublished = false 
        };
        _videoRepositoryMock.GetVideoInformationById(actualId).Returns(video);

        // Act
        await _sut.UpdateVideoInformation(userId, request);

        // Assert
        video.Name.Should().Be("Updated Name");
        video.Description.Should().Be("Updated Desc");
        video.IsPrivate.Should().BeTrue();
        video.IsPublished.Should().BeTrue();
        await _videoRepositoryMock.Received(1).Update(video);
    }

    #endregion

    #region BanVideo Tests
    
    [Fact]
    public async Task BanVideo_WhenAlreadyBanned_ShouldThrowApiException400()
    {
        // Arrange
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";
        var video = new VideoService.Models.Video { Name = "Some video",VideoId = actualId, IsBanned = true };
        
        _videoRepositoryMock.FindById(actualId).Returns(video);

        // Act
        Func<Task> act = async () => await _sut.BanVideo(prefixedVideoId);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 400 && e.Message == "Process action error");
    }
    
    #endregion
    
    #region IncrementVideoView Tests

    [Fact]
    public async Task IncrementVideoView_WhenPlatformIsYouTube_ShouldReturnEarly()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var youtubeVideoId = "yt_dQw4w9WgXcQ";

        // Act
        await _sut.IncrementVideoView(requestUserId, youtubeVideoId);

        // Assert
        await _videoRepositoryMock.DidNotReceive().FindUserView(Arg.Any<Guid>(), Arg.Any<Guid>());
    }

    [Fact]
    public async Task IncrementVideoView_WhenUserAlreadyViewed_ShouldReturnEarly()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        _videoRepositoryMock.FindUserView(requestUserId, actualId).Returns(true);

        // Act
        await _sut.IncrementVideoView(requestUserId, prefixedVideoId);

        // Assert
        await _videoRepositoryMock.DidNotReceive().AddUserView(Arg.Any<VideoService.Models.UserView>());
    }

    [Fact]
    public async Task IncrementVideoView_WhenNewView_ShouldAddViewAndUpdateCount()
    {
        // Arrange
        var requestUserId = Guid.NewGuid();
        var actualId = Guid.NewGuid();
        var prefixedVideoId = $"wb_{actualId}";

        var video = new VideoService.Models.Video { Name = "Some video",VideoId = actualId, Views = 5 };

        _videoRepositoryMock.FindUserView(requestUserId, actualId).Returns(false);
        _videoRepositoryMock.FindById(actualId).Returns(video);
        _videoRepositoryMock.CountUserView(actualId).Returns(6);

        // Act
        await _sut.IncrementVideoView(requestUserId, prefixedVideoId);

        // Assert
        await _videoRepositoryMock.Received(1).AddUserView(Arg.Is<VideoService.Models.UserView>(uv => 
            uv.UserId == requestUserId && uv.VideoId == actualId));
        
        video.Views.Should().Be(6);
        await _videoRepositoryMock.Received(1).Update(video);
    }

    #endregion
}