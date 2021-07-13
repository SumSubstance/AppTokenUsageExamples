import com.fasterxml.jackson.databind.ObjectMapper;
import model.Applicant;
import model.DocSet;
import model.DocType;
import model.HttpMethod;
import model.IdDocSetType;
import model.Metadata;
import model.RequiredIdDocs;
import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;
import okhttp3.ResponseBody;
import okio.Buffer;
import org.apache.commons.codec.binary.Hex;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.time.Instant;
import java.util.Arrays;
import java.util.Collections;
import java.util.UUID;

public class AppTokenJavaExample {
    // The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
    private static final String SUMSUB_SECRET_KEY = "YOUR_SUMSUB_SECRET_KEY"; // Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
    private static final String SUMSUB_APP_TOKEN = "YOUR_SUMSUB_APP_TOKEN"; // Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
    private static final String SUMSUB_TEST_BASE_URL = "https://test-api.sumsub.com"; // Please don't forget to change when switching to production

    private static final ObjectMapper objectMapper = new ObjectMapper();

    public static void main(String[] args) throws IOException, InvalidKeyException, NoSuchAlgorithmException {
        // The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases

        // Such actions are presented below:
        // 1) Creating an applicant
        // 2) Adding a document to the applicant
        // 3) Getting applicant status
        // 4) Getting access token

        String applicantId = createApplicant();
        System.out.println("The applicant was successfully created: " + applicantId);

        String imageId = addDocument(applicantId, new File(AppTokenJavaExample.class.getResource("/images/sumsub-logo.png").getFile()));
        System.out.println("Identifier of the added document: " + imageId);

        String applicantStatusStr = getApplicantStatus(applicantId);
        System.out.println("Applicant status (json string): " + applicantStatusStr);

        String accessTokenStr = getAccessToken(applicantId);
        System.out.println("Access token (json string): " + accessTokenStr);
    }

    public static String createApplicant() throws IOException, NoSuchAlgorithmException, InvalidKeyException {
        // https://developers.sumsub.com/api-reference/#creating-an-applicant

        Applicant applicant = new Applicant(
                UUID.randomUUID().toString(),
                new RequiredIdDocs(Arrays.asList(identityDocSet, selfieDocSet))
        );

        Response response = sendPost(
                "/resources/applicants?levelName=basic-kyc-level",
                RequestBody.create(
                        objectMapper.writeValueAsString(applicant),
                        MediaType.parse("application/json; charset=utf-8")));

        ResponseBody responseBody = response.body();

        return responseBody != null ? objectMapper.readValue(responseBody.string(), Applicant.class).getId() : null;
    }

    public static String addDocument(String applicantId, File doc) throws NoSuchAlgorithmException, InvalidKeyException, IOException {
        // https://developers.sumsub.com/api-reference/#adding-an-id-document

        RequestBody requestBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("metadata", objectMapper.writeValueAsString(new Metadata(DocType.PASSPORT, "DEU")))
                .addFormDataPart("content", doc.getName(), RequestBody.create(doc, MediaType.parse("image/*")))
                .build();

        Response response = sendPost("/resources/applicants/" + applicantId + "/info/idDoc", requestBody);
        return response.headers().get("X-Image-Id");
    }

    public static String getApplicantStatus(String applicantId) throws NoSuchAlgorithmException, InvalidKeyException, IOException {
        // https://developers.sumsub.com/api-reference/#getting-applicant-status-api

        Response response = sendGet("/resources/applicants/" + applicantId + "/requiredIdDocsStatus");

        ResponseBody responseBody = response.body();
        return responseBody != null ? responseBody.string() : null;
    }

    public static String getAccessToken(String applicantId) throws NoSuchAlgorithmException, InvalidKeyException, IOException {
        // https://developers.sumsub.com/api-reference/#access-tokens-for-sdks

        Response response = sendPost("/resources/accessTokens?userId=" + applicantId, RequestBody.create(new byte[0], null));

        ResponseBody responseBody = response.body();
        return responseBody != null ? responseBody.string() : null;
    }

    private static Response sendPost(String url, RequestBody requestBody) throws IOException, InvalidKeyException, NoSuchAlgorithmException {
        long ts = Instant.now().getEpochSecond();

        Request request = new Request.Builder()
                .url(SUMSUB_TEST_BASE_URL + url)
                .header("X-App-Token", SUMSUB_APP_TOKEN)
                .header("X-App-Access-Sig", createSignature(ts, HttpMethod.POST, url, requestBodyToBytes(requestBody)))
                .header("X-App-Access-Ts", String.valueOf(ts))
                .post(requestBody)
                .build();

        Response response = new OkHttpClient().newCall(request).execute();

        if (response.code() != 200 && response.code() != 201) {
            // https://developers.sumsub.com/api-reference/#errors
            // If an unsuccessful answer is received, please log the value of the "correlationId" parameter.
            // Then perhaps you should throw the exception. (depends on the logic of your code)
        }
        return response;
    }

    private static Response sendGet(String url) throws IOException, InvalidKeyException, NoSuchAlgorithmException {
        long ts = Instant.now().getEpochSecond();

        Request request = new Request.Builder()
                .url(SUMSUB_TEST_BASE_URL + url)
                .header("X-App-Token", SUMSUB_APP_TOKEN)
                .header("X-App-Access-Sig", createSignature(ts, HttpMethod.GET, url, null))
                .header("X-App-Access-Ts", String.valueOf(ts))
                .get()
                .build();

        Response response = new OkHttpClient().newCall(request).execute();

        if (response.code() != 200 && response.code() != 201) {
            // https://developers.sumsub.com/api-reference/#errors
            // If an unsuccessful answer is received, please log the value of the "correlationId" parameter.
            // Then perhaps you should throw the exception. (depends on the logic of your code)
        }
        return response;
    }

    private static String createSignature(long ts, HttpMethod httpMethod, String path, byte[] body) throws NoSuchAlgorithmException, InvalidKeyException {
        Mac hmacSha256 = Mac.getInstance("HmacSHA256");
        hmacSha256.init(new SecretKeySpec(SUMSUB_SECRET_KEY.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
        hmacSha256.update((ts + httpMethod.name() + path).getBytes(StandardCharsets.UTF_8));
        byte[] bytes = body == null ? hmacSha256.doFinal() : hmacSha256.doFinal(body);
        return Hex.encodeHexString(bytes);
    }

    public static byte[] requestBodyToBytes(RequestBody requestBody) throws IOException {
        Buffer buffer = new Buffer();
        requestBody.writeTo(buffer);
        return buffer.readByteArray();
    }

}



