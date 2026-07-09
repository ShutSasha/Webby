using System.Linq.Expressions;
using AutoMapper;
using FluentAssertions;
using Grpc.Core;
using NSubstitute;
using NSubstitute.ExceptionExtensions;
using Webby.UserService.Clients;
using Webby.UserService.Dtos.Complaint;
using Webby.UserService.Helpers.Exception;
using Webby.UserService.Interfaces.Repository;
using Webby.UserService.Interfaces.Service;
using Webby.UserService.Models.Enums;
using Webby.UserService.Services;

namespace Webby.FunctionalTests.Tests.Complaint;

public class ComplaintServiceTests
{
    private readonly IComplaintRepository _complaintRepositoryMock;
    private readonly IUserService _userServiceMock;
    private readonly IMapper _mapperMock;
    private readonly VideoGrpcService.VideoGrpcServiceClient _videoGrpcClientMock;

    private readonly ComplaintService _sut;

    public ComplaintServiceTests()
    {
        _complaintRepositoryMock = Substitute.For<IComplaintRepository>();
        _userServiceMock = Substitute.For<IUserService>();
        _mapperMock = Substitute.For<IMapper>();
        _videoGrpcClientMock = Substitute.For<VideoGrpcService.VideoGrpcServiceClient>();

        _sut = new ComplaintService(
            _complaintRepositoryMock,
            _userServiceMock,
            _mapperMock,
            _videoGrpcClientMock
        );
    }

    #region CreateComplaint Tests

    [Fact]
    public async Task CreateComplaint_WhenAuthorNotFound_ShouldThrowApiException404()
    {
        // Arrange
        var authorId = Guid.NewGuid();
        var sourceId = Guid.NewGuid();
        var request = new CreateComplaintRequest { TargetType = ComplaintTargetType.User, TargetId = sourceId.ToString(), ReasonType = "Test Reason" };

        _userServiceMock.GetById(authorId).Returns((UserService.Models.User)null);

        // Act
        Func<Task> act = async () => await _sut.CreateComplaint(authorId, request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 404 && e.Message == "Create complaint error");
    }

    [Fact]
    public async Task CreateComplaint_WhenTargetIsUserAndSelfComplaint_ShouldThrowApiException400()
    {
        // Arrange
        var authorId = Guid.NewGuid();
        var request = new CreateComplaintRequest 
        { 
            TargetType = ComplaintTargetType.User,
            TargetId = authorId.ToString(),
            ReasonType = "Test Reason"
        };

        _userServiceMock.GetById(authorId).Returns(new Webby.UserService.Models.User());

        // Act
        Func<Task> act = async () => await _sut.CreateComplaint(authorId, request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 400 && e.Message == "Create complaint error");
    }

    [Fact]
    public async Task CreateComplaint_WhenTargetIsVideoAndValid_ShouldAddComplaint()
    {
        // Arrange
        var authorId = Guid.NewGuid();
        var videoGuid = Guid.NewGuid();
        var prefixedVideoId = $"wb_{videoGuid}";
        
        var request = new CreateComplaintRequest 
        { 
            TargetType = ComplaintTargetType.Video,
            TargetId = prefixedVideoId,
            ReasonType = "Spam",
            AdditionalInfo = "Test info"
        };

        _userServiceMock.GetById(authorId).Returns(new Webby.UserService.Models.User());

        // Правильно мокаем успешный gRPC ответ
        var grpcResponse = new Google.Protobuf.WellKnownTypes.BoolValue { Value = true };
        var asyncCall = new AsyncUnaryCall<Google.Protobuf.WellKnownTypes.BoolValue>(
            Task.FromResult(grpcResponse), Task.FromResult(new Metadata()), () => Status.DefaultSuccess, () => new Metadata(), () => { });
        
        _videoGrpcClientMock.CheckVideoExistsAsync(Arg.Any<CheckVideoExistRequest>()).Returns(asyncCall);

        // Act
        await _sut.CreateComplaint(authorId, request);

        // Assert
        await _complaintRepositoryMock.Received(1).Add(Arg.Is<Webby.UserService.Models.Complaint>(c => 
            c.AuthorId == authorId &&
            c.TargetId == videoGuid &&
            c.TargetType == ComplaintTargetType.Video &&
            c.ReasonType == "Spam"
        ));
    }

    [Fact]
    public async Task CreateComplaint_WhenTargetIsVideoAndServiceUnavailable_ShouldThrowApiException503()
    {
        // Arrange
        var authorId = Guid.NewGuid();
        var videoGuid = Guid.NewGuid();
        var request = new CreateComplaintRequest 
        { 
            TargetType = ComplaintTargetType.Video,
            TargetId = videoGuid.ToString(),
            ReasonType = "Spam"
        };

        _userServiceMock.GetById(authorId).Returns(new Webby.UserService.Models.User());
        
        _videoGrpcClientMock.CheckVideoExistsAsync(Arg.Any<CheckVideoExistRequest>())
            .Throws(new RpcException(new Status(StatusCode.Unavailable, "Service is offline")));

        // Act
        Func<Task> act = async () => await _sut.CreateComplaint(authorId, request);

        // Assert
        await act.Should().ThrowAsync<ApiException>()
            .Where(e => e.StatusCode == 503 && e.Message == "Create complaint error");
    }

    #endregion

    #region GetUserComplaints Tests

    [Fact]
    public async Task GetUserComplaints_ShouldReturnMappedList()
    {
        // Arrange
        var userId = Guid.NewGuid();
        var complaints = new List<UserService.Models.Complaint> 
        { 
            new() { ComplaintId = Guid.NewGuid(), TargetId = userId } 
        };
        
        var dtoList = new List<ComplaintDto> { new() };

        _complaintRepositoryMock.GetByPredicate(Arg.Any<Expression<Func<Webby.UserService.Models.Complaint, bool>>>())
            .Returns(complaints.AsEnumerable()); 

        _mapperMock.Map<List<ComplaintDto>>(Arg.Any<List<Webby.UserService.Models.Complaint>>()).Returns(dtoList);

        // Act
        var result = await _sut.GetUserComplaints(userId);

        // Assert
        result.Should().BeEquivalentTo(dtoList);
    }

    #endregion
}