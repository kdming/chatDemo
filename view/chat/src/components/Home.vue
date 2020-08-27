<template>
  <div id="main">
    <p style="font-size:15px;">
      简易聊天室 项目地址:
      <a href="https://github.com/kdming/chat.git">https://github.com/kdming/chat.git</a>
    </p>
    <el-container>
      <el-header height="40px">
        <div id="header">
          <el-row :gutter="6" type="flex" justify="start" align="middle">
            <el-col :span="6">
              <div class="grid-content bg-purple">
                <el-tag>我的昵称:</el-tag>
                <span style="color:blue;margin-left:5px;font-size:15px;font-weight:bold;">{{user}}</span>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="grid-content bg-purple">
                <el-tag>在线人数:</el-tag>
                <span
                  style="color:red;margin-left:5px;font-size:15px;font-weight:bold;"
                >{{onelineNum}}</span>
              </div>
            </el-col>
          </el-row>
        </div>
      </el-header>
      <el-main>
        <div class="talk_show" id="words">
          <!-- 根据vue对象中的数组，遍历出对应的标签。 -->
          <div v-for="i in msgList" :class="i.type=='1'?'atalk':'btalk'" :key="i.words">
            <p class="talk-time">{{i.date}}</p>
            <span>{{ i.user }}：{{ i.msg }}</span>
          </div>
        </div>
      </el-main>
      <el-footer>
        <div id="send">
          <el-input
            type="textarea"
            placeholder="请输入内容,回车键发送"
            v-model="inputStr"
            maxlength="30"
            @keydown.native="sendMsg($event)"
            show-word-limit
          ></el-input>
        </div>
      </el-footer>
    </el-container>

    <!-- 用户名填写区域 -->
    <el-dialog
      title="登录"
      :visible.sync="userInputShow"
      width="30%"
      :before-close="handleClose"
      :show-close="false"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
    >
      <el-input v-model="user" placeholder="请填写用户名"></el-input>
      <!-- <span slot="footer" class="dialog-footer">
        <el-button type="primary"  @click="userInputShow = false">确 定</el-button>
      </span>-->
      <el-button
        type="primary"
        @click="checkUserInput"
        style="margin-top:40px;width:100px"
        size="medium"
      >确 定</el-button>
    </el-dialog>
  </div>
</template>

<script>
Date.prototype.format = function(fmt) {
  var o = {
    "M+": this.getMonth() + 1, //月份
    "d+": this.getDate(), //日
    "h+": this.getHours(), //小时
    "m+": this.getMinutes(), //分
    "s+": this.getSeconds(), //秒
    "q+": Math.floor((this.getMonth() + 3) / 3), //季度
    S: this.getMilliseconds() //毫秒
  };
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(
      RegExp.$1,
      (this.getFullYear() + "").substr(4 - RegExp.$1.length)
    );
  }
  for (var k in o) {
    if (new RegExp("(" + k + ")").test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length)
      );
    }
  }
  return fmt;
};
export default {
  name: "Home",
  data() {
    return {
      user: "",
      onelineNum: 0,
      userInputShow: false,
      msgList: [
        {
          user: "admin",
          msg: "Hi 快来聊天吧",
          date: new Date().format("yyyy-MM-dd hh:mm"),
          type: "1"
        }
      ],
      inputStr: "",
      wsuri: "ws://" + document.location.host + "/ws?key=", // ws wss
      lockReconnect: true, // 连接失败不进行重连
      maxReconnect: 5, // 最大重连次数，若连接失败
      socket: null
    };
  },
  mounted() {
    // const nameKey = "userName-009-008xxxx";
    // if (!this.user) {
    //   localStorage.removeItem(nameKey);
    //   this.user = localStorage.getItem(nameKey);
    // }
    if (this.user === "" || !this.user) {
      this.userInputShow = true;
    } else {
      // localStorage.setItem(nameKey, this.user);
      this.initWebSocket();
    }
  },
  methods: {
    reconnect() {
      console.log("尝试重连");
      if (this.lockReconnect || this.maxReconnect <= 0) {
        return;
      }
      setTimeout(() => {
        // this.maxReconnect-- // 不做限制 连不上一直重连
        this.initWebSocket();
      }, 60 * 1000);
    },
    initWebSocket() {
      try {
        if ("WebSocket" in window) {
          this.socket = new WebSocket(this.wsuri + this.user);
        } else {
          console.log("您的浏览器不支持websocket");
        }
        console.log("连接成功");
        this.socket.onopen = this.websocketonopen;
        this.socket.onerror = this.websocketonerror;
        this.socket.onmessage = this.websocketonmessage;
        this.socket.onclose = this.websocketclose;
        this.onelineNum += 1;
      } catch (e) {
        this.reconnect();
      }
    },
    websocketonopen() {
      console.log("WebSocket连接成功", this.socket.readyState);
    },
    websocketonerror(e) {
      console.log("WebSocket连接发生错误", e);
      alert("连接失败");
      this.reconnect();
    },
    websocketonmessage(e) {
      let data = JSON.parse(e.data);
      this.onelineNum = data.onlineNum;
      // 判断消息类型
      if (data.type === "2") {
        this.$notify({
          title: "新用户上线通知",
          message: data.msg,
          type: "success"
        });
        return;
      }
      const msg = {
        user: data.user,
        msg: data.msg,
        date: data.date,
        type: "2"
      };
      this.msgList.push(msg);
      this.scrollToBottom();
    },
    websocketclose(e) {
      console.log("connection closed (" + e.code + ")");
      alert("连接断开，请重新连接");
      this.onelineNum = 0;
      this.reconnect();
    },
    websocketsend() {
      const msg = {
        user: this.user,
        msg: this.inputStr,
        date: new Date().format("yyyy-MM-dd hh:mm"),
        type: "1"
      };
      this.msgList.push(msg);
      this.socket.send(this.inputStr);
      this.scrollToBottom();
    },
    destroyed() {
      this.socket.close();
    },
    handleClose() {
      this.userInputShow = false;
    },
    sendMsg(event) {
      if (event.keyCode == 13) {
        if (!event.metaKey) {
          event.preventDefault();
          this.websocketsend();
          this.inputStr = "";
        } else {
          this.inputStr = this.inputStr + "\n";
        }
      }
    },
    checkUserInput() {
      if (this.user === "" || !this.user) {
        alert("用户名不能为空");
        return;
      }
      this.userInputShow = false;
      this.initWebSocket();
    },
    scrollToBottom() {
      this.$nextTick(() => {
        var container = this.$el.querySelector("#words");
        container.scrollTop = container.scrollHeight;
      });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#main {
  width: 100%;
  padding-top: 0px;
  padding-bottom: 80px;
}
#header {
  width: 60%;
  height: 40px;
  margin: 0px auto 0;
  /* margin: 0% 20% 0% 20%; */
  /* vertical-align: center; */
  padding-top: 10px;
  /* border: solid 2px rgb(134, 197, 216); */
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12), 0 0 6px rgba(0, 0, 0, 0.04);
}

.el-main {
  margin: 0px !important;
}

.talk_show {
  width: 60%;
  height: 450px;
  /* border: 2px solid #e0878d; */
  background: #fff;
  margin: 10px auto 0;
  padding-top: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.12), 0 0 6px rgba(0, 0, 0, 0.04);
  /* box-shadow: 0 2px 4px rgba(29, 28, 28, 0.12), 0 0 6px rgba(0, 0, 0, 0.04); */
  overflow: auto;
}

#send {
  width: 60%;
  margin: 0% 20% 0% 20%;
  height: auto;
  text-align: center;
  border: solid 1px rgb(123, 200, 236);
}

/* 聊天内容A css */
.atalk {
  margin: 10px;
  text-align: left;
}
.atalk span {
  display: inline-block;
  /* background: #0181cc; */
  background: rgb(226, 249, 253);
  border-radius: 10px;
  color: rgb(90, 84, 84);
  padding: 5px 10px;
}

/* 聊天内容B css */
.btalk {
  margin: 10px;
  text-align: right;
}
.btalk span {
  display: inline-block;
  background: #cf9377;
  border-radius: 10px;
  color: rgb(250, 248, 248);
  padding: 5px 10px;
}

.talk-time {
  font-size: 4px;
  color: #0181cc;
  margin-left: 4px;
}
</style>
