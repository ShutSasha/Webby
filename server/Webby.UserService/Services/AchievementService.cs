using Microsoft.AspNetCore.Mvc;
using Webby.UserService.Dtos.Achievement;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models;

namespace Webby.UserService.Services;

public class AchievementService : IAchievementService
{
    private readonly IAchievementRepository _achievementRepository;
    private readonly IStorageService _storageService;

    public AchievementService(
        IAchievementRepository achievementRepository,
        IStorageService storageService)
    {
        _achievementRepository = achievementRepository;
        _storageService = storageService;
    }

    public async Task<Achievement> CreateAchievement(CreateAchievementRequest request)
    {
        var achievementExist = (await _achievementRepository
                .GetByPredicate(a => a.Title == request.Title || a.Description == request.Description))
            .FirstOrDefault();

        if (achievementExist != null)
        {
            throw new ApiException(
                "Create achievement error",
                400,
                "Achievement with this data already exist");
        }

        if (request.File.Length == 0)
        {
            throw new ApiException(
                "Create achievement error",
                400,
                "Achievement file icon is empty");
        }

        await using var fileStream = request.File.OpenReadStream();
        var achievementId = Guid.NewGuid();

        var achievementIconPath = await _storageService.UploadFileAsync(
            achievementId,
            "achievements_data",
            request.File.FileName,
            fileStream,
            request.File.ContentType);

        var achievement = new Achievement
        {
            AchievementId = achievementId,
            Code = request.Code,
            Description = request.Description,
            IconUrl = achievementIconPath,
            Title = request.Title
        };

        await _achievementRepository.Add(achievement);

        return achievement;
    }

    public async Task DeleteAchievement(Guid achievementId)
    {
        var achievement = await _achievementRepository.FindById(achievementId);

        if (achievement == null)
        {
            throw new ApiException(
                "Delete achievement error",
                404,
                "Achievement not found");
        }

        await _achievementRepository.DeleteAsync(achievement.AchievementId);
        
        if (!string.IsNullOrEmpty(achievement.IconUrl))
        {
            await _storageService.DeleteFileAsync(achievement.IconUrl);
        }
    }

    public async Task<Achievement> UpdateAchievement(UpdateAchievementRequest request)
    {
        var achievement = (await _achievementRepository
                .GetByPredicate(a => a.AchievementId == request.AchievementId))
            .FirstOrDefault();

        if (achievement == null)
        {
            throw new ApiException(
                "Update achievement error",
                404,
                "Achievement not found");
        }

        achievement.Title = request.Title;
        achievement.Description = request.Description;
        achievement.Code = request.Code;
        
        if (request.File.Length != 0)
        {
            await using var fileStream = request.File.OpenReadStream();

            var newIconPath = await _storageService.UploadFileAsync(
                achievement.AchievementId,
                "achievements_data",
                request.File.FileName,
                fileStream,
                request.File.ContentType);

            await _storageService.DeleteFileAsync(achievement.IconUrl);
            
            achievement.IconUrl = newIconPath;
        }

        await _achievementRepository.Update(achievement);

        return achievement;
    }

    public async Task<List<Achievement>> GetAchievements()
    {
        var achievements = await _achievementRepository.GetAll();
        return achievements.ToList();
    }
}