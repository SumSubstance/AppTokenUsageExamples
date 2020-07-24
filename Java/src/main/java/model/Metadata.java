package model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class Metadata {
    // https://developers.sumsub.com/api-reference/#request-metadata-body-part-fields
    private DocType idDocType;
    private String country;

    public Metadata() {
    }

    public Metadata(DocType idDocType, String country) {
        this.idDocType = idDocType;
        this.country = country;
    }

    public DocType getIdDocType() {
        return idDocType;
    }

    public void setIdDocType(DocType idDocType) {
        this.idDocType = idDocType;
    }

    public String getCountry() {
        return country;
    }

    public void setCountry(String country) {
        this.country = country;
    }
}
