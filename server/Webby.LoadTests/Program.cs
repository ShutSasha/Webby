using NBomber.Contracts.Stats;
using NBomber.CSharp;
using Webby.LoadTests.Scenarios;

var httpClient = new HttpClient();
httpClient.BaseAddress = new Uri("http://localhost:5000");

var httpSearchScenario = HttpVideoScenarios.HttpSearchWithText(httpClient);

NBomberRunner
   .RegisterScenarios(httpSearchScenario)
   .WithReportFormats(ReportFormat.Html, ReportFormat.Md) 
   .Run();