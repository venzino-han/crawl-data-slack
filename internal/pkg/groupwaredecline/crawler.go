package groupwaredecline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/moonsub-kim/crawl-data-slack/internal/pkg/crawler"
	"go.uber.org/zap"
)

type Crawler struct {
	logger       *zap.Logger
	ctx          context.Context
	eventBuilder eventBuilder
	id           string
	pw           string
	masters      []string
}

func (c Crawler) GetCrawlerName() string { return "groupware" }
func (c Crawler) GetJobName() string     { return "declined_payments" }

func (c Crawler) Crawl() ([]crawler.Event, error) {
	var jsonBody string
	var dtos []DTO

	startDate := time.Now().Add(time.Hour*9).AddDate(0, 0, -3).Format("2006-01-02") // UTCNOW -> KST -> 3days ago -> formating
	c.logger.Info("startDate", zap.String("startDate", startDate))
	err := chromedp.Run(
		c.ctx,
		chromedp.Navigate("https://gr.buzzvil.com/gw/uat/uia/egovLoginUsr.do"),

		// 로그인페이지: 로그인
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`document.getElementById('userId').value = '%s'; document.getElementById('userPw').value = '%s'; actionLogin();`,
				c.id,
				c.pw,
			),
			nil,
		),
		chromedp.Sleep(time.Second*2),

		// master 반려함
		chromedp.Navigate(`https://gr.buzzvil.com/eap/admin/eadoc/EADocMngList.do?menu_no=1705020000`),

		// 전자 결재 - 반려함: 일자 조정 버튼
		chromedp.EvaluateAsDevTools(
			fmt.Sprintf(
				`document.getElementById('from_date').value = '%s';
				document.getElementById('txt_doc_no').value = '지출결의서';
				document.querySelector('#ddlDocSts').value='100'; // 반려상태
				BindPuddGrid();`,
				startDate,
			),
			nil,
		),
		chromedp.Sleep(time.Second*2),

		// 문서 파싱
		chromedp.Evaluate(
			`
			function map_object(arr) {
				const indexMap = {2: "uid", 3: "doc_name", 4: "request_date", 6: "drafter", 7: "status"};
				var keys = Object.keys(indexMap);
				var obj = {};

				for (var i = 0; i < keys.length; i++) {
					k = indexMap[keys[i]];
					obj[k] = arr[keys[i]];
				}

				return obj;
			}

			function crawl() {
				var trs = document.body.querySelectorAll('div.grid-content > table > tbody > tr');
				var records = [];
				if (trs.length == 0) {
					return "[]" 			// ignore empty search results
				}

				for (var i = 0; i < trs.length; i++) {
					var arr = [];
					var tds = trs[i].getElementsByTagName('td');
					for (var j = 0; j < tds.length; j++) {
						arr.push(tds[j].innerText);
					}
					console.log(arr)
					records.push(map_object(arr));
				}

				return JSON.stringify(records);
			}
			crawl();
			`,
			&jsonBody,
		),
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonBody), &dtos)
	if err != nil {
		return nil, err
	}

	if len(dtos) == 1 && dtos[0].isEmpty() {
		c.logger.Warn("no data parsed")
		return []crawler.Event{}, nil
	}

	c.logger.Info("dto", zap.Any("dto", dtos))
	events, err := c.eventBuilder.buildEvents(dtos, c.GetCrawlerName(), c.GetJobName(), c.masters)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewCrawler(logger *zap.Logger, chromectx context.Context, id string, pw string, masters []string) *Crawler {
	return &Crawler{
		logger:  logger,
		ctx:     chromectx,
		id:      id,
		pw:      pw,
		masters: masters,
	}
}
