package paginator

import (
	"fmt"
	"gohub/pkg/config"
	"gohub/pkg/logger"
	"math"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Page struct {
	CurrentPage int    // 当前页
	PerPage     int    // 每页条数
	TotalPage   int    // 总页数
	TotalCount  int64  // 总条数
	NextPageUrl string // 下一页的链接
	PrevPageUrl string // 上一页的链接
}

type Paginator struct {
	BaseURL     string // 用以拼接 URL
	PerPage     int    // 每页条数
	CurrentPage int    // 当前页
	Offset      int    // 数据库读取数据时 Offset 的值
	TotalCount  int64  // 总页数
	TotalPage   int    // 总页数 = TotalCount/PerPage
	Sort        string // 排序规则
	Order       string // 排序顺序

	query *gorm.DB     // db query 句柄
	ctx   *gin.Context // gin context，方便调用
}

// Paginate 分页
// c —— gin.context 用来获取分页的 URL 参数
// db —— GORM 查询句柄，用以查询数据集和获取数据总数
// baseURL —— 用以分页链接
// data —— 模型数组，传址获取数据
// PerPage —— 每页条数，优先从 url 参数里取，否则使用 perPage 的值
// 用法:
//         query := database.DB.Model(Topic{}).Where("category_id = ?", cid)
//      var topics []Topic
//         page := paginator.Paginate(
//             c,
//             query,
//             &topics,
//             app.APIURL(database.TableName(&Topic{})),
//             perPage,
//         )
func Paginate(c *gin.Context, db *gorm.DB, data interface{}, baseURL string, perPage int) Page {
	// 初始化 Paginator 实例
	p := &Paginator{
		query: db,
		ctx:   c,
	}
	p.initProperties(perPage, baseURL)

	// 查询数据库
	err := p.query.Preload(clause.Associations). // 读取关联
							Order(p.Sort + " " + p.Order).
							Limit(p.PerPage).
							Offset(p.Offset).
							Find(data).
							Error
	if err != nil {
		logger.LogIf(err)
		return Page{}
	}

	return Page{
		CurrentPage: p.CurrentPage,
		PerPage:     p.PerPage,
		TotalPage:   p.TotalPage,
		TotalCount:  p.TotalCount,
		NextPageUrl: p.getNextPageURL(),
		PrevPageUrl: p.getPrevPageURL(),
	}
}

// initProperties 初始化分页必须用到的属性，基于这些属性查询数据库
func (p *Paginator) initProperties(perPage int, baseURL string) {
	p.BaseURL = p.formatBaseURL(baseURL)
	p.PerPage = p.getPerPage(perPage)

	// 排序参数（控制器中以验证过这些参数，可放心使用）
	p.Order = p.ctx.DefaultQuery(config.Get("page.url_query_order"), "asc")
	p.Sort = p.ctx.DefaultQuery(config.Get("page.url_query_sort"), "id")

	p.TotalCount = p.getTotalCount()
	p.TotalPage = p.getTotalPage()
	p.CurrentPage = p.getCurrentPage()
	p.Offset = (p.CurrentPage - 1) * p.PerPage
}

// getPerPage 获取
func (p *Paginator) getPerPage(perPage int) int {
	queryPerPage := p.ctx.Query(config.Get("page.url_query_per_page"))
	if len(queryPerPage) > 0 {
		perPage = cast.ToInt(queryPerPage)
	}

	// 没有传参，使用默认
	if perPage <= 0 {
		perPage = config.GetInt("page.perpage")
	}
	return perPage
}

// getCurrentPage 返回当前页码
func (p *Paginator) getCurrentPage() int {
	page := cast.ToInt(p.ctx.Query(config.Get("page.url_query_page")))
	if page <= 0 {
		page = 1
	}
	// TotalPage 等于 0 ，意味着数据不够分页
	if p.TotalCount == 0 {
		return 0
	}
	// 请求页数大于总页数，返回总页数
	if page > p.TotalPage {
		return p.TotalPage
	}
	return page
}

// getTotalCount 返回的是数据库里的条数
func (p *Paginator) getTotalCount() int64 {
	var count int64

	if err := p.query.Count(&count).Error; err != nil {
		return 0
	}
	return count
}

// getTotalPage 计算总页数
func (p *Paginator) getTotalPage() int {
	if p.TotalCount == 0 {
		return 0
	}
	nums := int64(math.Ceil(float64(p.TotalCount) / float64(p.PerPage)))
	if nums == 0 {
		nums = 1
	}
	return int(nums)
}

// formatBaseURL 兼容 URL 带与不带 `?` 的情况
func (p *Paginator) formatBaseURL(baseURL string) string {
	if strings.Contains(baseURL, "?") {
		baseURL += "&"
	} else {
		baseURL += "?"
	}

	return baseURL + config.Get("page.url_query_page") + "="
}

// getPageLink 拼接分页链接
func (p *Paginator) getPageLink(page int) string {
	return fmt.Sprintf("%v%v&%s=%s&%s=%s&%s=%v",
		p.BaseURL,
		page,
		config.Get("page.url_query_sort"),
		p.Sort,
		config.Get("page.url_query_order"),
		p.Order,
		config.Get("page.url_query_per_page"),
		p.PerPage,
	)
}

// getNextPageURL 返回下一页的链接
func (p *Paginator) getNextPageURL() string {
	if p.TotalPage > p.CurrentPage {
		return p.getPageLink(p.CurrentPage + 1)
	}
	return ""
}

// getPrevPageURL 返回下一页的链接
func (p *Paginator) getPrevPageURL() string {
	if p.CurrentPage == 1 || p.CurrentPage > p.TotalPage {
		return ""
	}
	return p.getPageLink(p.CurrentPage - 1)
}
