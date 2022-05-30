package httplib

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/imaimang/mgo/idlib"
)

func DigestAuth(ip string, port int, authURL string, method string, userName, userPwd string) (string, error) {
	url := "http://" + ip + ":" + strconv.Itoa(port) + "/" + authURL
	authorizationHeader := ""
	client := new(http.Client)
	request, err := http.NewRequest("GET", url, nil)
	if err == nil {
		var response *http.Response
		response, err = client.Do(request)
		if err == nil {
			if response.StatusCode == http.StatusUnauthorized {
				authenticates := response.Header.Values("Www-Authenticate")
				var realm string
				//var domain string
				var qop string
				var nonce string
				var opaque string
				//var algorithm string
				//var stale bool
				for _, authenticate := range authenticates {
					if strings.Index(authenticate, "Digest") == 0 {
						authenticate = authenticate[7:]
						items := strings.Split(authenticate, ",")
						for _, item := range items {
							itemsKeyValue := strings.Split(item, "=")
							key := strings.TrimSpace(strings.ToUpper(itemsKeyValue[0]))
							value := strings.TrimSpace(strings.Trim(itemsKeyValue[1], "\""))
							switch key {
							case "REALM":
								realm = value
							case "DOMAIN":
								//domain = value
							case "QOP":
								qop = value
							case "NONCE":
								nonce = value
							case "OPAQUE":
								opaque = value
							case "ALGORITHM":
								//algorithm = value

							case "STALE":
								//stale, _ = strconv.ParseBool(value)
							}

						}
						break
					}
				}
				nc := time.Now().Format("060201150405")
				cnonce := idlib.CreateMD5()
				ha1 := idlib.Encrypt32MD5(userName + ":" + realm + ":" + userPwd)
				ha2 := idlib.Encrypt32MD5(method + ":" + authURL)
				encryptResult := idlib.Encrypt32MD5(ha1 + ":" + nonce + ":" + nc + ":" + cnonce + ":" + qop + ":" + ha2)
				authorizationHeader = "Digest  username=\"" + userName + "\", realm=\"" + realm + "\", nonce=\"" + nonce + "\", uri=\"" + authURL + "\", " + "qop=\"" + qop + "\", nc=\"" + nc + "\", cnonce=\"" + cnonce + "\", response=\"" + encryptResult + "\", opaque=\"" + opaque + "\""
				request, err = http.NewRequest(method, "http://"+ip+":"+strconv.Itoa(port)+"/"+authURL, nil)
				if err == nil {
					request.Header.Add("Authorization", authorizationHeader)
					response, err = client.Do(request)
					if err == nil {
						if response.StatusCode != http.StatusOK {
							err = errors.New(response.Status)
						}
					}
				}
			}
		}
	}
	return authorizationHeader, err
}
