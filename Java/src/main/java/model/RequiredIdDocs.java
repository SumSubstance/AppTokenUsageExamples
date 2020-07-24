package model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

import java.util.List;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class RequiredIdDocs {
    private List<DocSet> docSets;

    public RequiredIdDocs() {
    }

    public RequiredIdDocs(List<DocSet> docSets) {
        this.docSets = docSets;
    }

    public List<DocSet> getDocSets() {
        return docSets;
    }

    public void setDocSets(List<DocSet> docSets) {
        this.docSets = docSets;
    }
}
