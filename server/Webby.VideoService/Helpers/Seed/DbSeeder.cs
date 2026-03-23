using Webby.VideoService.Interfaces.Repositories;
using Webby.VideoService.Interfaces.Services;
using Webby.VideoService.Models;

namespace Webby.VideoService.Helpers.Seed;

public class DbSeeder
{
   private readonly IVideoRepository _repository;
   private readonly IStorageService _storageService;
   private readonly Random _rng = new();
   public static readonly string[] VideoTitles = 
    {
        "Top 10 Hidden Gems in Southern Europe",
        "How to Make the Perfect Homemade Pizza",
        "Morning Routine for Maximum Productivity",
        "The Mystery of the Abandoned Mansion",
        "Urban Exploration: Exploring New York at Night",
        "Walking Through Tokyo: 4K ASMR Journey",
        "5 Easy Exercises to Improve Your Posture",
        "My Experience Living Off the Grid for a Week",
        "The Science of Why We Dream",
        "15 Mind-Blowing Space Facts You Didn't Know",
        "Street Food Tour in Bangkok: Spicy Delights",
        "How to Start a Small Garden on Your Balcony",
        "Learning a New Language in 3 Months: My Tips",
        "The Evolution of Cinema: From Silent to 8K",
        "10 Life Hacks for Frequent Travelers",
        "Relaxing Rainy Night in a Cozy Cabin",
        "Budget Travel: Iceland on $50 a Day",
        "The Secret History of the Great Pyramids",
        "Unboxing the Latest Tech Gadgets of 2026",
        "Meditation for Beginners: Finding Inner Peace",
        "Why Coffee is Actually Good for Your Brain",
        "The Art of Minimalism: Decluttering My Life",
        "Exploring the World's Most Dangerous Roads",
        "How to Capture Professional Photos with Your Phone",
        "DIY Home Decor: Transforming My Bedroom",
        "The Future of Electric Cars: What to Expect",
        "A Weekend in London: Best Places to Visit",
        "Cooking Challenge: 3 Meals with Only $10",
        "The Most Isolated Tribes on Earth",
        "Philosophy 101: Stoicism in Modern Life"
    };

    public DbSeeder(IVideoRepository repository, IStorageService storageService)
    {
        _repository = repository;
        _storageService = storageService;
    }

    public async Task Seed(Guid authorId, string sourceFolderPath, int totalCount)
    {
        var sourceVideos = Directory.GetFiles(sourceFolderPath, "*.mp4");
        var sourceImages = Directory.GetFiles(sourceFolderPath, "*.jpg");

        if (sourceVideos.Length == 0) throw new System.Exception("Source folder has no videos!");

        for (int i = 0; i < totalCount; i++)
        {
            var videoId = Guid.NewGuid();
            string sourceVideoPath = sourceVideos[i % sourceVideos.Length];
            string sourceImagePath = sourceImages.Length > 0 ? sourceImages[i % sourceImages.Length] : "";
            
            string videoUrl;
            using (var videoStream = File.OpenRead(sourceVideoPath))
            {
                videoUrl = await _storageService.UploadFileAsync(
                    videoId, 
                    "videos", 
                    Path.GetFileName(sourceVideoPath), 
                    videoStream, 
                    "video/mp4"
                );
            }

            string thumbUrl = "https://i.pinimg.com/736x/d0/33/dd/d033dddf7ee3da574ad482a145da802e.jpg";
            if (!string.IsNullOrEmpty(sourceImagePath))
            {
                using (var imageStream = File.OpenRead(sourceImagePath))
                {
                    thumbUrl = await _storageService.UploadFileAsync(
                        videoId, 
                        "previews", 
                        Path.GetFileName(sourceImagePath), 
                        imageStream, 
                        "image/jpeg"
                    );
                }
            }
            
            var video = new Models.Video
            {
                VideoId = videoId,
                UserId = authorId,
                Name = VideoTitles[i],
                Description = "Auto-generated seed content uploaded to S3.",
                VideoUrl = videoUrl,
                PreviewUrl = thumbUrl,
                IsPrivate = false,
                Views = _rng.Next(0, 5000),
                CreatedAt = DateTime.UtcNow,
            };

            await _repository.Add(video);
        }
    }
}