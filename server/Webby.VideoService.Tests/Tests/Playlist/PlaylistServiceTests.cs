using System.Linq.Expressions;
using AutoMapper;
using FluentAssertions;
using NSubstitute;
using UserService;
using Webby.VideoService.Constants;
using Webby.VideoService.Dtos.External;
using Webby.VideoService.Dtos.Playlist;
using Webby.VideoService.Dtos.Search;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Interfaces.Helpers;
using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;
using Webby.VideoService.Models.Enums;
using Webby.VideoService.Services;
using Xunit;

namespace Webby.VideoService.Tests.Tests.Playlist;

public class PlaylistServiceTests
{
   private readonly IPlaylistRepository _playlistRepositoryMock;
   private readonly IMapper _mapperMock;
   private readonly UserGrpcService.UserGrpcServiceClient _userClientMock;
   private readonly IVideoRepository _videoRepositoryMock;
   private readonly IYouTubeSearchService _youtubeSearchServiceMock;
   private readonly ITwitchSearchService _twitchSearchServiceMock;
   private readonly IExternalContentFetcher _externalContentFetcherMock;
   
   private readonly PlaylistService _sut;

   public PlaylistServiceTests()
   {
      _playlistRepositoryMock = Substitute.For<IPlaylistRepository>();
      _mapperMock = Substitute.For<IMapper>();
      _userClientMock = Substitute.For<UserGrpcService.UserGrpcServiceClient>();
      _videoRepositoryMock = Substitute.For<IVideoRepository>();
      _youtubeSearchServiceMock = Substitute.For<IYouTubeSearchService>();
      _twitchSearchServiceMock = Substitute.For<ITwitchSearchService>();
      _externalContentFetcherMock = Substitute.For<IExternalContentFetcher>();
      
      _sut = new PlaylistService(
         _playlistRepositoryMock,
         _mapperMock,
         _userClientMock,
         _videoRepositoryMock,
         _youtubeSearchServiceMock,
         _twitchSearchServiceMock,
         _externalContentFetcherMock
      );
   }
   
   [Fact]
   public async Task CreatePlaylist_ShouldCallRepositoryAdd_AndReturnMappedDto()
   {
      // Arrange
      var userId = Guid.NewGuid();
      var request = new CreatePlaylistRequest { Name = "My Favorite Anime OSTs", IsPrivate = false };
        
      var expectedDto = new PlaylistDto { Name = "My Favorite Anime OSTs", IsPrivate = false };
      _mapperMock.Map<PlaylistDto>(Arg.Any<Models.Playlist>()).Returns(expectedDto);

      // Act
      var result = await _sut.CreatePlaylist(userId, request);

      // Assert
      result.Should().BeEquivalentTo(expectedDto);
      
      await _playlistRepositoryMock.Received(1).Add(Arg.Is<Models.Playlist>(p => 
         p.Name == request.Name && 
         p.UserId == userId &&
         p.IsPrivate == request.IsPrivate));
   }
   
   [Fact]
   public async Task GetPlaylistById_WhenPlaylistNotFound_ShouldThrowApiException404()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      _playlistRepositoryMock.FindById(playlistId).Returns((Models.Playlist)null);

      // Act
      Func<Task> act = async () => await _sut.GetPlaylistById(playlistId);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 404 && e.Message == "Get playlist error");
   }

   [Fact]
   public async Task GetPlaylistById_WhenPlaylistIsPrivate_ShouldThrowApiException403()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var privatePlaylist = new Models.Playlist { PlaylistId = playlistId, IsPrivate = true };
      _playlistRepositoryMock.FindById(playlistId).Returns(privatePlaylist);

      // Act
      Func<Task> act = async () => await _sut.GetPlaylistById(playlistId);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 403 && e.Message == "Get playlist error");
   }

   [Fact]
   public async Task GetPlaylistById_WhenPlaylistExistsAndPublic_ShouldReturnPlaylist()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var publicPlaylist = new Models.Playlist { PlaylistId = playlistId, IsPrivate = false };
      _playlistRepositoryMock.FindById(playlistId).Returns(publicPlaylist);

      // Act
      var result = await _sut.GetPlaylistById(playlistId);

      // Assert
      result.Should().NotBeNull();
      result.PlaylistId.Should().Be(playlistId);
      result.IsPrivate.Should().BeFalse();
   }

   [Fact]
   public async Task UpdatePlaylist_WhenUserIsNotOwner_ShouldThrowApiException403()
   {
      // Arrange
      var requestUserId = Guid.NewGuid();
      var actualOwnerId = Guid.NewGuid();
      var playlistId = Guid.NewGuid();
      
      var request = new UpdatePlaylistRequest 
      { 
         PlaylistId = playlistId, 
         Name = "Hacked Name", 
         IsPrivate = false 
      };
      
      var existingPlaylist = new Models.Playlist 
      { 
         PlaylistId = playlistId, 
         UserId = actualOwnerId 
      };
      
      _playlistRepositoryMock.FindById(playlistId).Returns(existingPlaylist);
      
      // Act
      
      Func<Task> act = async () => await _sut.UpdatePlaylist(requestUserId, request);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 403 && 
                     e.Message == "Update playlist error");
      
      await _playlistRepositoryMock.DidNotReceive().Update(Arg.Any<Models.Playlist>());
   }

   [Fact]
   public async Task UpdatePlaylist_WhenPlaylistNotFound_ShouldReturn404()
   {
      // Arrange
      var playlistId  = Guid.NewGuid();
      var requestUserId = Guid.NewGuid();

      var request = new UpdatePlaylistRequest()
      {
         Name = "abob",
         IsPrivate = false,
         PlaylistId = playlistId
      };

      // Act
      var act = async () => await _sut.UpdatePlaylist(requestUserId,request);
      
      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 404 && e.Message == "Update playlist error");

      await _playlistRepositoryMock.DidNotReceive().Update(Arg.Any<Models.Playlist>());
   }

   [Fact]
   public async Task UpdatePlaylist_ShouldReturnPlaylist()
   {
      // Arrange
      var actualOwnerId = Guid.NewGuid();
      var playlistId = Guid.NewGuid();
      
      var request = new UpdatePlaylistRequest 
      { 
         PlaylistId = playlistId, 
         Name = "Hacked Name", 
         IsPrivate = false 
      };
      
      var existingPlaylist = new Models.Playlist 
      { 
         PlaylistId = playlistId, 
         UserId = actualOwnerId,
         IsPrivate = true,
         Name = "Aboba"
      };
      
      _playlistRepositoryMock.FindById(playlistId).Returns(existingPlaylist);
      
      //Act
      var result = await _sut.UpdatePlaylist(actualOwnerId, request);
      
      //Assert
      result.Should().NotBeNull();
      result.PlaylistId.Should().Be(playlistId);
      result.IsPrivate.Should().BeFalse();
      result.Name.Should().Be(request.Name);

   }
   
   [Fact]
   public async Task DeletePlaylist_WhenPlaylistNotFound_ShouldThrowApiException404()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var userId = Guid.NewGuid();
      _playlistRepositoryMock.FindById(playlistId).Returns((Models.Playlist)null);

      // Act
      Func<Task> act = async () => await _sut.DeletePlaylist(userId, playlistId);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 404 && e.Message == "Delete playlist error");
   }

   [Fact]
   public async Task DeletePlaylist_WhenUserIsNotOwner_ShouldThrowApiException403()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var requestUserId = Guid.NewGuid();
      var ownerId = Guid.NewGuid(); // Другой пользователь

      var playlist = new Models.Playlist { PlaylistId = playlistId, UserId = ownerId };
      _playlistRepositoryMock.FindById(playlistId).Returns(playlist);

      // Act
      Func<Task> act = async () => await _sut.DeletePlaylist(requestUserId, playlistId);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 403 && e.Message == "Delete playlist error");
    
      // Убеждаемся, что до удаления дело не дошло
      await _playlistRepositoryMock.DidNotReceive().DeleteAsync(Arg.Any<Guid>());
   }

   [Fact]
   public async Task DeletePlaylist_WhenValid_ShouldCallRepositoryDelete()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var userId = Guid.NewGuid();
      var playlist = new Models.Playlist { PlaylistId = playlistId, UserId = userId };
    
      _playlistRepositoryMock.FindById(playlistId).Returns(playlist);

      // Act
      await _sut.DeletePlaylist(userId, playlistId);

      // Assert
      await _playlistRepositoryMock.Received(1).DeleteAsync(playlistId);
   }
   
   [Fact]
   public async Task GetPlaylistInformation_WhenValid_ShouldReturnMappedResponse()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var requestedUserId = Guid.NewGuid();
    
      var playlist = new Models.Playlist 
      { 
         PlaylistId = playlistId, 
         Name = "Chainsaw Man Openings",
         PlaylistVideos = new List<PlaylistVideo>()
      };

      var expectedDto = new PlaylistDto { Name = "Chainsaw Man Openings" };

      _playlistRepositoryMock.FindByIdWithVideos(playlistId).Returns(playlist);
      _mapperMock.Map<PlaylistDto>(Arg.Any<Models.Playlist>()).Returns(expectedDto);
      _externalContentFetcherMock
         .FetchExternalContentAsync(Arg.Any<List<PlaylistVideo>>())
         .Returns(new ExternalContentData()); 

      // Act
      var result = await _sut.GetPlaylistInformation(playlistId, requestedUserId);

      // Assert
      result.Should().NotBeNull();
      result.Playlist.Name.Should().Be("Chainsaw Man Openings");
      result.HiddenVideosCount.Should().Be(0);
   }

   [Fact]
   public async Task GetPlaylistInformation_ShouldHidePrivateVideosFromOtherUsers()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var requestedUserId = Guid.NewGuid();
      var otherUserId = Guid.NewGuid();
 
      var playlist = new Models.Playlist 
      { 
         PlaylistId = playlistId,
         PlaylistVideos = new List<PlaylistVideo>
         {
            new PlaylistVideo 
            { 
               Platform = SystemPlatforms.Webby, 
               Video = new Video(isPrivate: true, userId: otherUserId, name: "new video #1")
            }
         }
      };

      _playlistRepositoryMock.FindByIdWithVideos(playlistId).Returns(playlist);

      _externalContentFetcherMock
         .FetchExternalContentAsync(Arg.Any<List<PlaylistVideo>>())
         .Returns(new ExternalContentData()); 

      // Act
      var result = await _sut.GetPlaylistInformation(playlistId, requestedUserId);

      // Assert
      result.HiddenVideosCount.Should().Be(1);
      result.Playlist.CountOfVideos.Should().Be(0);
   }
   
   [Fact]
   public async Task AttachVideoToPlaylist_WhenVideoItemsEmpty_ShouldThrowApiException400()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var userId = Guid.NewGuid();
      var emptyList = new List<string>();

      // Act
      Func<Task> act = async () => await _sut.AttachVideoToPlaylist(playlistId, emptyList, userId);

      // Assert
      await act.Should().ThrowAsync<ApiException>()
         .Where(e => e.StatusCode == 400 && e.Message == "Attach video error");
   }

   [Fact]
   public async Task CheckIfVideoExistInPlaylist_ShouldReturnRepositoryResult()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var actualId = Guid.NewGuid().ToString();
      var videoId = PlatformPrefixesConstants.WebbyPrefix + actualId;
    
      _playlistRepositoryMock.CheckIsVideoAdded(actualId, playlistId).Returns(true);

      // Act
      var result = await _sut.CheckIfVideoExistInPlaylist(playlistId, videoId);

      // Assert
      result.Should().BeTrue();
      await _playlistRepositoryMock.Received(1).CheckIsVideoAdded(actualId, playlistId);
   }

   [Fact]
   public async Task AttachVideoToPlaylist_WhenNewLocalVideo_ShouldCallAddPlaylistVideos()
   {
      // Arrange
      var playlistId = Guid.NewGuid();
      var userId = Guid.NewGuid();
      var videoId = PlatformPrefixesConstants.WebbyPrefix + Guid.NewGuid();

      var requestItems = new List<string>() { videoId };

      var playlist = new Models.Playlist 
      { 
         PlaylistId = playlistId, 
         UserId = userId, 
         PlaylistVideos = new List<PlaylistVideo>() 
      };

      _playlistRepositoryMock
         .GetByPredicate(Arg.Any<Expression<Func<Models.Playlist, bool>>>())
         .Returns(new List<Models.Playlist> { playlist }.AsEnumerable());

      var itemsToMockReturn = new List<PlaylistVideo> 
      { 
         new PlaylistVideo { Platform = SystemPlatforms.Webby } 
      };

      _playlistRepositoryMock
         .GetPlaylistItemsDiffAsync(playlistId, MediaType.Video, Arg.Any<List<PlaylistVideo>>())
         .Returns((itemsToMockReturn, new List<PlaylistVideo>()));
   
      _videoRepositoryMock.CheckVideosCount(Arg.Any<List<Guid>>()).Returns(true);
      _videoRepositoryMock.CheckForbiddenVideos(Arg.Any<List<Guid>>(), userId).Returns(false);

      // Act
      await _sut.AttachVideoToPlaylist(playlistId, requestItems, userId);

      // Assert
      await _playlistRepositoryMock.Received(1).AddPlaylistVideos(Arg.Is<List<PlaylistVideo>>(list => 
         list.Count == 1 && list[0].Platform == SystemPlatforms.Webby));
   }
}