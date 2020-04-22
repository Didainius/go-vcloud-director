package govcd

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"time"
)

func getSamlRequestBody(user, pass, samlEntityIderence string) string {
	return `<s:Envelope xmlns:s="http://www.w3.org/2003/05/soap-envelope" xmlns:a="http://www.w3.org/2005/08/addressing" xmlns:u="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd">
	<s:Header>
	  <a:Action s:mustUnderstand="1">http://docs.oasis-open.org/ws-sx/ws-trust/200512/RST/Issue</a:Action>
	  <a:ReplyTo>
		<a:Address>http://www.w3.org/2005/08/addressing/anonymous</a:Address>
	  </a:ReplyTo>
	  <a:To s:mustUnderstand="1">https://win-60g606n0afg.test-forest.net/adfs/services/trust/13/usernamemixed</a:To>
	  <o:Security s:mustUnderstand="1" xmlns:o="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd">
		<u:Timestamp u:Id="_0">
		  <u:Created>` + time.Now().Format(time.RFC3339) + `</u:Created>
		  <u:Expires>` + time.Now().Add(1*time.Minute).Format(time.RFC3339) + `</u:Expires>
		</u:Timestamp>
		<o:UsernameToken>
		  <o:Username>` + user + `</o:Username>
		  <o:Password o:Type="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText">` + pass + `</o:Password>
		</o:UsernameToken>
	  </o:Security>
	</s:Header>
	<s:Body>
	  <trust:RequestSecurityToken xmlns:trust="http://docs.oasis-open.org/ws-sx/ws-trust/200512">
		<wsp:AppliesTo xmlns:wsp="http://schemas.xmlsoap.org/ws/2004/09/policy">
		  <a:samlEntityIderence>
			<a:Address>` + samlEntityIderence + `</a:Address>
		  </a:samlEntityIderence>
		</wsp:AppliesTo>
		<trust:KeySize>0</trust:KeySize>
		<trust:KeyType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Bearer</trust:KeyType>
		<i:RequestDisplayToken xml:lang="en" xmlns:i="http://schemas.xmlsoap.org/ws/2005/05/identity" />
		<trust:RequestType>http://docs.oasis-open.org/ws-sx/ws-trust/200512/Issue</trust:RequestType>
		<trust:TokenType>http://docs.oasis-open.org/wss/oasis-wss-saml-token-profile-1.1#SAMLV2.0</trust:TokenType>
	  </trust:RequestSecurityToken>
	</s:Body>
  </s:Envelope>
  `
}

func gzipAndBase64Encode(token string) (string, error) {
	var gzipBuffer bytes.Buffer
	gz := gzip.NewWriter(&gzipBuffer)
	if _, err := gz.Write([]byte(token)); err != nil {
		return "", fmt.Errorf("error writing to gzip buffer: %s", err)
	}
	if err := gz.Close(); err != nil {
		return "", fmt.Errorf("error closing gzip buffer: %s", err)
	}
	base64GzippedToken := base64.StdEncoding.EncodeToString(gzipBuffer.Bytes())

	return base64GzippedToken, nil
}
