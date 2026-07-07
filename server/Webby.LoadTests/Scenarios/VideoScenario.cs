using System;
using System.Net.Http;
using System.Net.Http.Json;
using NBomber.Contracts;
using NBomber.CSharp;
using Webby.VideoService.Dtos.Video;
using Webby.VideoService.Dtos.Video.Enums;

namespace Webby.LoadTests.Scenarios;

public static class HttpVideoScenarios
{
   public static ScenarioProps HttpSearchWithText(HttpClient httpClient)
   {
      return Scenario.Create("http_text_search", async context =>
         {
            var searchWord = DataFeeds.SearchWordsFeed.GetNextItem(context.ScenarioInfo);

            var requestBody = new SearchVideoOptions
            {
               SearchText = searchWord,
               Page = 1,
               PageSize = 20,
               SearchPlatform = SearchVideoPlatforms.Webby
            };

            var request = new HttpRequestMessage(HttpMethod.Get, "/api/videos/search")
            {
               Content = JsonContent.Create(requestBody)
            };

            try
            {
               var response = await httpClient.SendAsync(request);

               if (response.IsSuccessStatusCode)
               {
                  var size = response.Content.Headers.ContentLength ?? 0;
                  return Response.Ok(sizeBytes: (int)size);
               }

               var errorBody = await response.Content.ReadAsStringAsync();
               return Response.Fail(statusCode: ((int)response.StatusCode).ToString(), message: errorBody);
            }
            catch (Exception ex)
            {
               return Response.Fail(statusCode: "500", message: ex.Message);
            }
         })
         .WithLoadSimulations(
            Simulation.RampingConstant(copies: 20, during: TimeSpan.FromSeconds(10)),
            Simulation.KeepConstant(copies: 20, during: TimeSpan.FromMinutes(1))
         );
   }
}