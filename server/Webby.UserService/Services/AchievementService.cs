using AutoMapper;
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
    private readonly IMapper _mapper;
    public AchievementService(
        IAchievementRepository achievementRepository,
        IStorageService storageService, IMapper mapper)
    {
        _achievementRepository = achievementRepository;
        _storageService = storageService;
        _mapper = mapper;
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
            Description = request.Description,
            IconUrl = achievementIconPath,
            Title = request.Title,
            TargetValue = request.TargetValue,
            EventType = request.EventType
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
        achievement.EventType = request.EventType;
        achievement.TargetValue = request.TargetValue;
        
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
        => await _achievementRepository.GetAll();

    public async Task<Achievement?> FindById(Guid achievementId)
    {
        return await _achievementRepository.FindById(achievementId);
    }

    public async Task AddUserAchievement(Guid userId, Guid achievementId)
    {
        var userAchievementExist = await _achievementRepository
            .HasUserAchievement(userId, achievementId);

        if (userAchievementExist)
        {
            throw new ApiException("Add user achievement error", 400, "user achievement is already exist");
        }

        await _achievementRepository.AddUserAchievement(userId, achievementId);
    }

    public async Task<UserAchievement> GetUserAchievement(Guid userId, Guid achievementId)
    {
        var userAchievement = await _achievementRepository.GetUserAchievement(userId, achievementId);

        if (userAchievement == null)
        {
            throw new ApiException("Get user achievement error", 400, "This achievement has not been unlocked yet.");
        }

        return userAchievement;
    }

    public async Task UpdateUserAchievement(UserAchievement userAchievement) => 
        await _achievementRepository.UpdateUserAchievement(userAchievement);

    public async Task<List<ProfileAchievementDto>> GetPinnedAchievements(Guid userId)
    {
        var userAchievements = await _achievementRepository
            .GetUserPinnedAchievements(userId);

        return userAchievements
            .Select(ua => _mapper.Map<ProfileAchievementDto>(ua.Achievement))
            .ToList();
    }

    public async Task<GetUserAchievementsResponse> GetUserAchievementsBlock(Guid userId)
    {
        var getUserAchievementsResponse = new GetUserAchievementsResponse();
        
        var pinnedAchievements = await _achievementRepository
            .GetUserPinnedAchievements(userId);

        var achievementsWithUserStatus = await _achievementRepository
            .GetAchievementsWithUserStatus(userId);
        
        getUserAchievementsResponse.PinnedAchievements = pinnedAchievements.Select(
            ua => _mapper.Map<AchievementDto>(ua.Achievement)).ToList();

        getUserAchievementsResponse.Achievements = achievementsWithUserStatus;

        return getUserAchievementsResponse;
    }

    public async Task<int> GetPinnedAchievementsCount(Guid userId) 
        => await _achievementRepository.CountPinnedAchievements(userId);
}