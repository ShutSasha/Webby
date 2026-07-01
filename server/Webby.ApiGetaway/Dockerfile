FROM mcr.microsoft.com/dotnet/aspnet:9.0-alpine AS base
USER $APP_UID
WORKDIR /app
ENV ASPNETCORE_HTTP_PORTS=8000
EXPOSE 8000

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["Webby.ApiGetaway/Webby.ApiGetaway.csproj", "Webby.ApiGetaway/"]
RUN dotnet restore "Webby.ApiGetaway/Webby.ApiGetaway.csproj"
COPY . .
WORKDIR "/src/Webby.ApiGetaway"
RUN dotnet build "Webby.ApiGetaway.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "Webby.ApiGetaway.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --chown=app:app --from=publish /app/publish .
ENTRYPOINT ["dotnet", "Webby.ApiGetaway.dll"]