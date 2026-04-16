using Webby.VideoService.Dtos.Stream;
using Webby.VideoService.Dtos.Stream.Enums;
using Webby.VideoService.Helpers.Exception;
using Webby.VideoService.Helpers.Response;
using Webby.VideoService.Interfaces.Services;

namespace Webby.VideoService.Services;

public class StreamService : IStreamService
{
   private readonly TwitchSearchService _twitchSearchService;

   public StreamService(TwitchSearchService twitchSearchService)
   {
      _twitchSearchService = twitchSearchService;
   }
   
   public async Task<PagedResponse<StreamDto>> SearchStream(SearchStreamOptions searchOptions)
   {
      searchOptions.SearchPlatform ??= SearchStreamPlatforms.Twitch;

      return searchOptions.SearchPlatform switch
      {
         SearchStreamPlatforms.Twitch =>
            await _twitchSearchService.SearchAsync(searchOptions.SearchText, searchOptions.PageSize, searchOptions.Page,
               searchOptions.NextPageToken),

         _ => throw new ApiException("Search stream error", 400, "Invalid type of search stream platform")
      };
   }

   public async Task<StreamDto> GetStreamById(string streamerId, SearchStreamPlatforms? platforms)
   {
      platforms ??= SearchStreamPlatforms.Twitch;
      
      return platforms switch
      {
         SearchStreamPlatforms.Twitch =>
            await _twitchSearchService.FindById(streamerId),

         _ => throw new ApiException("Get stream error", 400, "Invalid type of search stream platform")
      };
   }
   
}