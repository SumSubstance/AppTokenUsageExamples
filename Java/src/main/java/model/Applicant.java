package model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class Applicant {
    // https://developers.sumsub.com/api-reference/#request-body
    private String id;
    private String externalUserId;
    private RequiredIdDocs requiredIdDocs;

    public Applicant() {
    }

    public Applicant(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getExternalUserId() {
        return externalUserId;
    }

    public void setExternalUserId(String externalUserId) {
        this.externalUserId = externalUserId;
    }

    public RequiredIdDocs getRequiredIdDocs() {
        return requiredIdDocs;
    }

    public void setRequiredIdDocs(RequiredIdDocs requiredIdDocs) {
        this.requiredIdDocs = requiredIdDocs;
    }
}
