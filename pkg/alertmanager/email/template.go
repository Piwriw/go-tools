package email

var SMTPDefaultPort = 587
var EmailTableTemplate = `
<div class="table_content">
        <span class="line"></span>
        <table>
          <tbody>
            <tr>
              <td class="name">告警项:</td>
              <td class="detail">%s</td>
            </tr>

            <tr>
              <td class="name">告警描述:</td>
              <td class="detail">%s</td>
            </tr>

            <tr>
              <td class="name">监控对象:</td>
              <td class="detail">%s</td>
            </tr>
            <tr>
                <td class="name">监控IP:</td>
                <td class="detail">%s</td>
              </tr>
            <tr>
              <td class="name">告警开始时间:</td>
              <td class="detail">%s</td>
            </tr>
            <tr>
              <td class="name">告警级别:</td>
              <td class="detail">%s</td>
            </tr>
            <tr>
              <td class="name">状态:</td>
              <td class="detail">%s</td>
            </tr>
            <tr>
              <td class="name">告警信息:</td>
              <td class="detail">%s</td>
            </tr>
            <tr>
              <td class="name">建议:</td>
              <td class="detail">%s</td>
            </tr>
          </tbody>
        </table>
      </div>
`
var EmailTemplateHeader = `
<div><font>
    </font>
</div>
<div><includetail><!--<![endif]--></includetail></div>
<div class="notification">
      <style>
        .notification {
          background-color: #f2f2f2;
        }
        table {
          border: 1px solid #f2f2f2;
          padding: 0.5% 2% 2% 1%;
          background-color: #ffffff;
          width: 98%;
          border-collapse: separate;
          border-spacing: 0 10px;
          vertical-align: middle;
        }
        .detail {
          color: #777777;
          background-color: #f2f2f2;
          width: auto;
          padding: 15px;
          margin: 5px;
        }
        .alert_title {
          font-size: large;
          vertical-align: middle;
        }
        .line {
          width: 0.3%;
          height: 60px;
          background-color: black;
          border: solid black 3px;
          vertical-align: middle;
          float: left;
        }
        .table_title {
          margin-left: 5px;
        }
        .table_content {
          margin-top: 5px;
        }
      </style>
      <div class="table_title">
        <svg
          t="1614670116148"
          class="icon"
          style="vertical-align: middle"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="2190"
          width="35"
          height="35"
        >
          <path
            d="M784 144H240c-88 0-160 72-160 160v416c0 88 72 160 160 160h544c88 0 160-72 160-160V304c0-88-72-160-160-160z m-544 64h544c35.2 0 65.6 19.2 81.6 46.4L552 508.8c-24 19.2-57.6 19.2-80 0L158.4 254.4c16-27.2 46.4-46.4 81.6-46.4z m544 608H240c-52.8 0-96-43.2-96-96V326.4L432 560c24 19.2 51.2 28.8 80 28.8 28.8 0 57.6-9.6 80-28.8l288-233.6V720c0 52.8-43.2 96-96 96z"
            p-id="2191"
            fill="#EA4D4E"
          ></path>
        </svg>
        <span class="alert_title">告警通知</span>
      </div>
`
