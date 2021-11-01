require 'json'
require 'rest-client'
require_relative 'model/applicant'
require_relative 'model/fixed_info'
require_relative 'model/metadata'
require_relative 'model/id_doc_types'
require_relative 'model/jsonable'
require 'securerandom'

# The description of the authorization method is available here: https://developers.sumsub.com/api-reference/#app-tokens
APP_TOKEN = 'YOUR_SUMSUB_APP_TOKEN'.freeze # Example: sbx:uY0CgwELmgUAEyl4hNWxLngb.0WSeQeiYny4WEqmAALEAiK2qTC96fBad
SECRET_KEY = 'YOUR_SUMSUB_SECRET_KEY'.freeze # Example: Hej2ch71kG2kTd1iIUDZFNsO5C1lh5Gq
# Please don't forget to change token and secret key values to production ones when switching to production

def request_env_url(resource)
  "https://api.sumsub.com/resources/#{resource}"
end

# https://developers.sumsub.com/api-reference/#creating-an-applicant

def create_applicant(external_used_id, lang, level_name)
  resources = "applicants?levelName=#{level_name}"
  body = Applicant.new(external_used_id, lang).serialize
  puts body
  header = signed_header(resources, body)

  response = RestClient.post(
    request_env_url(resources),
    body,
    header
  )
end

# https://developers.sumsub.com/api-reference/#adding-an-id-document

def upload_photo(applicant_id)
  resources = "applicants/#{applicant_id}/info/idDoc"
  bounds = '-------------------------12431412541'

  image = File.open('resources/images/sumsub-logo.png', 'rb') { |io| io.read }
  payload = "#{+'--'}#{bounds}\r\n"
  payload += "Content-Disposition: form-data; name=\"metadata\"\r\nContent-type: application/json; charset=utf-8\r\n\r\n"
  payload += Metadata.new(IdDocType::PASSPORT, 'GBR').serialize
  payload += "\r\n"
  payload += "--#{bounds}\r\n"
  payload += "Content-Disposition: form-data; name=\"content\"; filename=\"image.png\"\r\nContent-type: image/png; charset=utf-8\r\n\r\n"
  payload += "#{image}\r\n"
  payload += "--#{bounds}--"

  RestClient.post(request_env_url(resources),
                  payload,
                  signed_header(resources, payload, 'POST', "multipart/form-data; boundary=#{bounds}"))
end

# https://developers.sumsub.com/api-reference/#getting-applicant-status-api

def get_applicant_status(applicant_id)
  resources = "applicants/#{applicant_id}/requiredIdDocsStatus"
  RestClient.get request_env_url(resources), signed_header(resources, nil, 'GET')
end

# https://developers.sumsub.com/api-reference/#getting-applicant-data

def get_applicant_data(applicant_id)
  resources = "applicants/#{applicant_id}/one"

  response = RestClient.get request_env_url(resources), signed_header(resources, nil, 'GET')
end

# https://developers.sumsub.com/api-reference/#access-tokens-for-sdks

def generate_access_token(external_user_id, level_name = 'basic-kyc-level', ttl = 600)
  raise 'VIOLATION: Null id' if external_user_id.empty?

  # Send the request
  resources = "accessTokens?userId=#{external_user_id.to_s}&ttlInSecs=#{ttl.to_s}&levelName=#{level_name}"
  response = RestClient.post(request_env_url(resources), nil, signed_header(resources, nil, 'POST'))

  JSON.parse(response.body)
end

# https://developers.sumsub.com/api-reference/#app-tokens headers example

def signed_header(resource, body = nil, method = 'POST', content_type = 'application/json')
  epoch_time = Time.now.to_i
  access_signature = signed_message(epoch_time, resource, body, method)
  {
    "X-App-Token": APP_TOKEN.encode('UTF-8').to_s,
    "X-App-Access-Sig": access_signature.encode('UTF-8').to_s,
    "X-App-Access-Ts": epoch_time.to_s.encode('UTF-8').to_s,
    "Accept": 'application/json',
    "Content-Type": content_type.to_s
  }
end

# https://developers.sumsub.com/api-reference/#app-tokens

def signed_message(time, resource, body, method = 'POST')
  key = SECRET_KEY
  body_encoded = body.to_s if body
  data = "#{time.to_s}#{method}/resources/#{resource.to_s}#{body_encoded.to_s}"
  digest = OpenSSL::Digest.new('sha256')
  OpenSSL::HMAC.hexdigest(digest, key, data)
end

uuid = SecureRandom.uuid
level_name = 'basic-kyc-level'
lang = 'en'

# 1. create an applicant
response = JSON.parse(create_applicant(uuid, lang, level_name))
puts("applicant_id #{response['id']}")
applicant_id = response['id']

# 2. generate an access token
puts(generate_access_token(uuid, level_name))

# 3. upload a photo
puts(upload_photo(applicant_id))

# 4. getting applicant data
puts(get_applicant_data(applicant_id))

# 5. getting applicant status
puts(get_applicant_status(applicant_id))

