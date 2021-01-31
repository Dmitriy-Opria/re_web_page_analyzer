package api_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/api"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/model"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service"
	"github.com/Dmitriy-Opria/re_web_page_analyzer/internal/service/servicefakes"
	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Parsing Test", func() {

	var (
		staff   *service.StaffService
		fetcher *servicefakes.FakeFetcher
		parser  *service.ParserService
		router  http.Handler

		req = model.ParserRequest{
			URL: "https://www.w3schools.com/",
		}

		htmlPage = []byte(`<!DOCTYPE HTML>
			<html lang="en-US" xmlns="http://www.w3.org/1999/html" xmlns="http://www.w3.org/1999/html">
				<head>
					<title>W3Schools Online Web Tutorials</title>
				</head>
			<body>
			
			<h1>Header Level 1</h1>
			<h2>Header Level 2</h2>
			<h3>Header Level 3</h3>
			<h4>Header Level 4</h4>
			<h5>Header Level 5</h5>
			<h6>Header Level 6</h6>
			<div>
				<a href="/html/tryit.asp?filename=tryhtml_default">EXERCISES</a>
				<a href="/cert/default.asp" id="cert_navbtn">CERTIFICATES</a>
				<div >
					<a href="https://www.linkedin.com/company/w3schools.com/" title="W3Schools on LinkedIn">LinkedIn</a>
					<a href="https://www.instagram.com/w3schools.com_official/" title="W3Schools on Instagram">Instagram</a>
					<a href="https://www.facebook.com/w3schoolscom/" title="W3Schools on Facebook">Facebook</a>
				</div>
				<form>
					<label>Login</label>
					<label>Username</label>
					<input>Name</input>
					<label>Password</label>
					<input>Password</input>
					<button>Sign In</button>
				</form>
			</div>
			
			</body>
			</html>`)
	)
	JustBeforeEach(func() {
		staff = service.NewStaffService()
		fetcher = &servicefakes.FakeFetcher{}
		parser = service.NewParserService(fetcher, 1)

		router = api.NewHandler(staff, parser)
	})
	Describe("should parse html page and compare results", func() {
		JustBeforeEach(func() {
			doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(htmlPage))
			doc.Url = &url.URL{Host: "www.w3schools.com"}
			fetcher.FetchReturns(doc, nil)

			fetcher.IsAccessibleReturnsOnCall(0, &model.WorkerWrapper{Index: 0, Result: true}, nil)
			fetcher.IsAccessibleReturnsOnCall(1, &model.WorkerWrapper{Index: 1, Result: true}, nil)
			fetcher.IsAccessibleReturnsOnCall(2, &model.WorkerWrapper{Index: 2, Result: true}, nil)

			fetcher.IsAccessibleReturnsOnCall(3, &model.WorkerWrapper{Index: 0, Result: true}, nil)
			fetcher.IsAccessibleReturnsOnCall(4, &model.WorkerWrapper{Index: 1, Result: true}, nil)
			fetcher.IsAccessibleReturnsOnCall(5, &model.WorkerWrapper{Index: 2, Result: true}, nil)
		})
		It("should make parsing http request and check the results", func() {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(req)
			request, _ := http.NewRequest(http.MethodPost, "/api/v1/parsing/page/analyze", bytes.NewBuffer(reqBody))

			router.ServeHTTP(w, request)

			body, _ := ioutil.ReadAll(w.Result().Body)
			response := model.ParserResponse{}
			err := json.Unmarshal(body, &response)

			Expect(err).To(BeNil())

			Expect(response.Version).To(BeEquivalentTo("HTML5 and beyond"))
			Expect(response.Title).To(BeEquivalentTo("W3Schools Online Web Tutorials"))

			Expect(response.ListH1).To(BeEquivalentTo([]string{"Header Level 1"}))
			Expect(response.ListH2).To(BeEquivalentTo([]string{"Header Level 2"}))
			Expect(response.ListH3).To(BeEquivalentTo([]string{"Header Level 3"}))
			Expect(response.ListH4).To(BeEquivalentTo([]string{"Header Level 4"}))
			Expect(response.ListH5).To(BeEquivalentTo([]string{"Header Level 5"}))
			Expect(response.ListH6).To(BeEquivalentTo([]string{"Header Level 6"}))

			Expect(response.InternalLinks[0]).To(BeEquivalentTo(&model.Link{
				Name:       "EXERCISES",
				Url:        "http://www.w3schools.com/htmltryit.asp?filename=tryhtml_default",
				Accessible: true,
			}))
			Expect(response.InternalLinks[1]).To(BeEquivalentTo(&model.Link{
				Name:       "CERTIFICATES",
				Url:        "http://www.w3schools.com/certdefault.asp",
				Accessible: true,
			}))

			Expect(response.ExternalLinks[0]).To(BeEquivalentTo(&model.Link{
				Name:       "LinkedIn",
				Url:        "https://www.linkedin.com/company/w3schools.com/",
				Accessible: true,
			}))
			Expect(response.ExternalLinks[1]).To(BeEquivalentTo(&model.Link{
				Name:       "Instagram",
				Url:        "https://www.instagram.com/w3schools.com_official/",
				Accessible: true,
			}))
			Expect(response.ExternalLinks[2]).To(BeEquivalentTo(&model.Link{
				Name:       "Facebook",
				Url:        "https://www.facebook.com/w3schoolscom/",
				Accessible: true,
			}))

			Expect(response.Login).To(BeEquivalentTo(true))
		})
	})
})
