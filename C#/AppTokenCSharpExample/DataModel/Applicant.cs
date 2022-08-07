namespace AppTokenCSharpExample.DataModel;

public class Applicant
{
    public string id { get; set; }
    public string createdAt { get; set; }
    public string clientId { get; set; }
    public string inspectionId { get; set; }
    public string externalUserId { get; set; }
    public Fixedinfo fixedInfo { get; set; }
    public string email { get; set; }
    public string phone { get; set; }
    public Requirediddocs requiredIdDocs { get; set; }
    public Review review { get; set; }
    public string type { get; set; }
}