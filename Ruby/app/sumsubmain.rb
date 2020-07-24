#Ruby 2.6.6p146
require 'json'
require 'rest-client'

class Sumsub
  class << self

    # The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
    APP_TOKEN = "Some app token, that has to be generated in our dashboard" # Example: tst:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
    SECRET_KEY = "Some secret key, that has to be generated in our dashboard" # Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq


    # Please don't forget to change when switching to production
    def request_env_url(resource)
      'https://test-api.sumsub.com/resources/' + resource
    end

    # https://developers.sumsub.com/api-reference/#creating-an-applicant
    def applicantRequest(used_id, email, firstName, lastName)
      resources = 'applicants'
      body = {
          "externalUserId":"#{used_id}",
          "info":{
              "email":"#{email}",
              "firstName":"#{firstName}",
              "lastName":"#{lastName}"
          },
          "lang":"it",
          "requiredIdDocs": {
              "docSets":[
                  {
                      "idDocSetType":"IDENTITY",
                      "types":[
                          "PASSPORT",
                          "ID_CARD",
                          "DRIVERS",
                          "RESIDENCE_PERMIT"]
                  }
              ]
          }
      }.to_json
      header = signed_header(resources, body)

      response = RestClient.post(
          request_env_url(resources),
          body,
          header
      )
    end

    # https://developers.sumsub.com/api-reference/#adding-an-id-document
    def uploadPhoto(applicantId)
      resources = "applicants/" + applicantId + "/info/idDoc"
      bounds = "-------------------------12431412541"

      image = File.open('resources/images/sumsub-logo.png', 'rb') { |io| io.read }
      payload = + '--' + bounds + "\r\n"
      payload += "Content-Disposition: form-data; name=\"metadata\"" + "\r\nContent-type: application/json; charset=utf-8\r\n\r\n"
      payload += {
          "idDocType":"DRIVERS",
          "country":"USA"
      }.to_json
      payload += "\r\n"
      payload += '--' + bounds + "\r\n"
      payload += "Content-Disposition: form-data; name=\"content\"; filename=\"image.png\"\r\nContent-type: image/png; charset=utf-8\r\n\r\n"
      payload += image + "\r\n"
      payload += '--' + bounds + "--"

      RestClient.post(request_env_url(resources),
                      payload,
                      signed_header(resources, payload, "POST", "multipart/form-data; boundary="+bounds))

    end

    # https://developers.sumsub.com/api-reference/#getting-applicant-status-api
    def getApplicantStatus(applicantId)
      resources = "applicants/#{applicantId}/requiredIdDocsStatus"
      RestClient.get request_env_url(resources), signed_header(resources, nil, 'GET')
    end

    # https://developers.sumsub.com/api-reference/#getting-applicant-data
    def getApplicantData(applicantId)
      resources = "applicants/#{applicantId}/one"

      response = RestClient.get request_env_url(resources), signed_header(resources, nil, 'GET')
    end

    # https://developers.sumsub.com/api-reference/#access-tokens-for-sdks WEBSDK example
    def access_token(user_id, ttl=600)
      raise 'VIOLATION: Null id' if user_id.empty?
      # Send the request
      resources = 'accessTokens?userId='+ user_id.to_s + '&ttlInSecs='+ttl.to_s
      puts request_env_url(resources)
      response = RestClient.post( request_env_url(resources), nil, signed_header(resources, nil, 'POST'))

      JSON.parse(response.body)
    end

    # https://developers.sumsub.com/api-reference/#app-tokens headers example
    def signed_header(resource, body=nil, method='POST', contenttype="application/json")
      epoch_time = Time.now.to_i
      puts "Epoch: "+epoch_time.to_s
      access_signature = signed_message(epoch_time, resource, body, method)
      {
          "X-App-Token": "#{APP_TOKEN.encode("UTF-8")}",
          "X-App-Access-Sig": "#{access_signature.encode("UTF-8")}",
          "X-App-Access-Ts": "#{epoch_time.to_s.encode("UTF-8")}",
          "Accept": "application/json",
          "Content-Type": "#{contenttype}"
      }
    end

    # https://developers.sumsub.com/api-reference/#app-tokens signing message
    def signed_message(time, resource, body, method="POST")
      key = SECRET_KEY
      body_encoded = body.to_s if body
      data = time.to_s + method + '/resources/'+resource.to_s+ body_encoded.to_s
      digest = OpenSSL::Digest.new('sha256')
      OpenSSL::HMAC.hexdigest(digest, key, data)
    end





  end

  # The description of the flow can be found here: https://developers.sumsub.com/api-flow/#api-integration-phases
  #
  # Such actions are presented below:
  #
  # 1) Creating an applicant
  # 2) Adding a document to the applicant
  # 3) Getting applicant status
  # 4) Getting applicant data
  # 5) Initializing WEBSDK

  applicantRequest(4071505, 'support@sumsub.com', 'Clercqer', 'LastName')
  puts uploadPhoto("Some applicantId here") #For example: 5f000dd5aee05c701a7e8874
  getApplicantStatus("Some applicantId here") #For example: 5f000dd5aee05c701a7e8874
  getApplicantData("Some applicantId here") #For example: 5f000dd5aee05c701a7e8874
  access_token("4071505")

end