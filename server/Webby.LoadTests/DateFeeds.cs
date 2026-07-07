using NBomber.Contracts;
using NBomber.CSharp;
using NBomber.Data;
using NBomber.Data.CSharp;

namespace Webby.LoadTests;

public static class DataFeeds
{
   private static readonly string[] SearchKeywords = 
   {
      "miku", "test", "video", "aboba", "vlog", 
      "street", "morning", "top", "how", "webby"
   };
   
   public static IDataFeed<string> SearchWordsFeed = DataFeed.Random(SearchKeywords);
    
   public static IDataFeed<string> UserIdsFeed = DataFeed.Random(
      Enumerable.Range(0, 500).Select(_ => Guid.NewGuid().ToString()).ToArray());
}