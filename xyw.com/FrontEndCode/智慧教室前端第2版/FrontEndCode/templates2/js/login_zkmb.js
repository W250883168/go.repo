$(document).ready(function() {

	$("#username,#userpwd").change(function() {
		$("#showMeg").html("");
	});
	$("#Loginbutton").keydown(function(e) {
		var curKey = e.which;
		if(curKey == 13) {
			$("#Loginbutton").click();
			return false;
		}
	});
	$("#Loginbutton").click(function(e) {
		var Loginuser = $.trim($("#username").val());
		var Loginpwd = $.trim($("#userpwd").val());
		if(Loginuser == "" || Loginuser == "undefined") {
			$("#showMeg").html("<b style='color:#fff;'>请输入账号 </b>");
			return false;
		} else if(Loginpwd == "" || Loginpwd == "undefined") {
			$("#showMeg").html("<b style='color:#fff;'>请输入密码 </b>");
			return false;
		} else {
			//   显示加载中
			$(".icon-loading").show();
			var postdata = {
				"Loginuser": Loginuser,
				"Loginpwd": Loginpwd,
				"Os":"Controlpanel"
			};
			$.ajax({
				dataType: "JSON",
				url: "/login",
				type: "POST",
				//contentType: "application/json; charset=utf-8",
				data: JSON.stringify(postdata),
				success: function(data) {
					$(".icon-loading").hide();
					if(data.Rcode == "1000") {
						//localStorage.clear();
						localStorage.setItem("LoginUser", JSON.stringify(data.Result))
							//$.widows.location.href="index.html";
						window.location.href = "/web2/html/#/app/zkmb";
					} else {
						$("#showMeg").html("<b style='color:#FF6434;'>" + data.Reason + " </b>");
					}
				}
			});
		}
		
		return false;
	});
});