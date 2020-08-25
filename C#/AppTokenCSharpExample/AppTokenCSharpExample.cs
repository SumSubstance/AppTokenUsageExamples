using System;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.Http;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

namespace AppTokenCSharpExample
{
    internal class AppTokenCSharpExample
    {
        // The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
        private static readonly string SUMSUB_SECRET_KEY = "SUMSUB_SECRET_KEY"; // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
        private static readonly string SUMSUB_APP_TOKEN = "SUMSUB_APP_TOKEN";  // Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
        private static readonly string SUMSUB_TEST_BASE_URL = "https://test-api.sumsub.com";  // Please don't forget to change when switching to production

        private static void Main(string[] args)
        {
            // The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases

            // Such actions are presented below:
            // 1) Creating an applicant
            // 2) Adding a document to the applicant
            // 3) Getting applicant status
            // 4) Getting access token

            // Create an applicant
            string applicantId = CreateApplicant().Result.id;

            // Add a document to the applicant
            var addDocumentResult = AddDocument(applicantId).Result;
            Console.WriteLine("Add Document Result: " + ContentToString(addDocumentResult.Content));

            // Get Applicant Status
            var getApplicantResult = GetApplicantStatus(applicantId).Result;
            Console.WriteLine("Applicant status (json string): " + ContentToString(getApplicantResult.Content));

            // Get access token
            var accessTokenResult = GetAccessToken(applicantId).Result;
            Console.WriteLine("Access token Result: " + ContentToString(accessTokenResult.Content));

            // Important: please keep this line as async tasks that end unexpectedly will close console window before showing the error.
            Console.ReadLine();
        }

        // https://developers.sumsub.com/api-reference/#getting-applicant-status-sdk
        public static async Task<Applicant> CreateApplicant()
        {
            Console.WriteLine("Creating an applicant...");

            var body = new
            {
                externalUserId = $"USER_{DateTimeOffset.UtcNow.ToUnixTimeSeconds()}",
                requiredIdDocs = new
                {
                    docSets = new[]
                    {
                        new {idDocSetType = "IDENTITY", types = new[] {"PASSPORT", "ID_CARD", "DRIVERS"}},
                        new {idDocSetType = "SELFIE", types = new[] {"SELFIE"}}
                    }
                }
            };

            // Create the request body
            var requestBody = new HttpRequestMessage(HttpMethod.Post, SUMSUB_TEST_BASE_URL)
            {
                Content = new StringContent(JsonConvert.SerializeObject(body), Encoding.UTF8, "application/json")
            };

            // Get the response
            var response = await SendPost("/resources/applicants", requestBody);
            var applicant = JsonConvert.DeserializeObject<Applicant>(ContentToString(response.Content));

            Console.WriteLine(response.IsSuccessStatusCode
                ? $"The applicant was successfully created: {applicant.id}"
                : $"ERROR: {ContentToString(response.Content)}");

            return applicant;
        }

        // https://developers.sumsub.com/api-reference/#adding-an-id-document
        public static async Task<HttpResponseMessage> AddDocument(string applicantId)
        {
            Console.WriteLine("Adding document to the applicant...");

            // metadata object
            var metaData = new
            {
                idDocType = "PASSPORT",
                country = "GBR"
            };

            using (var formContent = new MultipartFormDataContent())
            {
                // Add metadata json object
                formContent.Add(new StringContent(JsonConvert.SerializeObject(metaData)), "\"metadata\"");

                // Add binary content
                var binaryImage = File.ReadAllBytes("../../resources/sumsub-logo.png");
                formContent.Add(new StreamContent(new MemoryStream(binaryImage)), "content", "sumsub-logo.png");

                // Request body
                var requestBody = new HttpRequestMessage(HttpMethod.Post, SUMSUB_TEST_BASE_URL)
                {
                    Content = formContent
                };

                var response = await SendPost($"/resources/applicants/{applicantId}/info/idDoc", requestBody);

                Console.WriteLine(response.IsSuccessStatusCode
                    ? $"Document was successfully added"
                    : $"ERROR: {ContentToString(response.Content)}");

                return response;
            }
        }

        // https://developers.sumsub.com/api-reference/#getting-applicant-status-api
        public static async Task<HttpResponseMessage> GetApplicantStatus(string applicantId)
        {
            Console.WriteLine("Getting the applicant status...");

            var response = await SendGet($"/resources/applicants/{applicantId}/requiredIdDocsStatus");
            return response;
        }

        public static async Task<HttpResponseMessage> GetAccessToken(string applicantId)
        {
            var response = await SendPost($"/resources/accessTokens?userId={applicantId}", new HttpRequestMessage());
            return response;
        }

        private static async Task<HttpResponseMessage> SendPost(string url, HttpRequestMessage requestBody)
        {

            var ts = DateTimeOffset.UtcNow.ToUnixTimeSeconds();
            var signature = CreateSignature(ts, HttpMethod.Post, url, RequestBodyToBytes(requestBody));

            ServicePointManager.SecurityProtocol = SecurityProtocolType.Tls12;
            var client = new HttpClient
            {
                BaseAddress = new Uri(SUMSUB_TEST_BASE_URL)
            };
            client.DefaultRequestHeaders.Add("X-App-Token", SUMSUB_APP_TOKEN);
            client.DefaultRequestHeaders.Add("X-App-Access-Sig", signature);
            client.DefaultRequestHeaders.Add("X-App-Access-Ts", ts.ToString());

            var response = await client.PostAsync(url, requestBody.Content);

            if (!response.IsSuccessStatusCode)
            {
                // https://developers.sumsub.com/api-reference/#errors
                // If an unsuccessful answer is received, please log the value of the "correlationId" parameter.
                // Then perhaps you should throw the exception. (depends on the logic of your code)
            }

            // debug
            //var debugInfo = response.Content.ReadAsStringAsync().Result;
            return response;
        }

        private static async Task<HttpResponseMessage> SendGet(string url)
        {
            long ts = DateTimeOffset.UtcNow.ToUnixTimeSeconds();

            ServicePointManager.SecurityProtocol = SecurityProtocolType.Tls12;
            var client = new HttpClient
            {
                BaseAddress = new Uri(SUMSUB_TEST_BASE_URL)
            };
            client.DefaultRequestHeaders.Add("X-App-Token", SUMSUB_APP_TOKEN);
            client.DefaultRequestHeaders.Add("X-App-Access-Sig", CreateSignature(ts, HttpMethod.Get, url, null));
            client.DefaultRequestHeaders.Add("X-App-Access-Ts", ts.ToString());

            var response = await client.GetAsync(url);

            if (!response.IsSuccessStatusCode)
            {
                // https://developers.sumsub.com/api-reference/#errors
                // If an unsuccessful answer is received, please log the value of the "correlationId" parameter.
                // Then perhaps you should throw the exception. (depends on the logic of your code)
            }

            return response;
        }

        private static string CreateSignature(long ts, HttpMethod httpMethod, string path, byte[] body)
        {
            Console.WriteLine("Creating a signature for the request...");

            var hmac256 = new HMACSHA256(Encoding.ASCII.GetBytes(SUMSUB_SECRET_KEY));

            byte[] byteArray = Encoding.ASCII.GetBytes(ts + httpMethod.Method + path);

            if (body != null)
            {
                // concat arrays: add body to byteArray
                var s = new MemoryStream();
                s.Write(byteArray, 0, byteArray.Length);
                s.Write(body, 0, body.Length);
                byteArray = s.ToArray();
            }

            var result = hmac256.ComputeHash(
                new MemoryStream(byteArray)).Aggregate("", (s, e) => s + String.Format("{0:x2}", e), s => s);

            return result;
        }

        private static string ContentToString(HttpContent httpContent)
        {
            return httpContent == null ? "" : httpContent.ReadAsStringAsync().Result;
        }

        private static byte[] RequestBodyToBytes(HttpRequestMessage requestBody)
        {
            return requestBody.Content == null ? 
                new byte[] { } : requestBody.Content.ReadAsByteArrayAsync().Result;
        }
    }
}
