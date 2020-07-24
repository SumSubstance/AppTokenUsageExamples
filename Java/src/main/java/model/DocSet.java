package model;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonInclude;

import java.util.List;

@JsonInclude(JsonInclude.Include.NON_NULL)
@JsonIgnoreProperties(ignoreUnknown = true)
public class DocSet {
    private IdDocSetType idDocSetType;
    private List<DocType> types;

    public DocSet() {
    }

    public DocSet(IdDocSetType idDocSetType, List<DocType> types) {
        this.idDocSetType = idDocSetType;
        this.types = types;
    }

    public IdDocSetType getIdDocSetType() {
        return idDocSetType;
    }

    public void setIdDocSetType(IdDocSetType idDocSetType) {
        this.idDocSetType = idDocSetType;
    }

    public List<DocType> getTypes() {
        return types;
    }

    public void setTypes(List<DocType> types) {
        this.types = types;
    }
}
