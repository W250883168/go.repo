$(document).ready(function() {
	$("#username").focus();
	$("#login").on("submit", function() {
		//console.log(1111)
		var Loginuser = $.trim($("#username").val());
		var Loginpwd = $.trim($("#userpwd").val());
		if(Loginuser == "" || Loginuser == "undefined") {
			$("#showMeg").html("<i class='icon-err'></i><span style='color:#FF9292;font-size:14px;'>请输入账号 </span>");
			return false;
		} else if(Loginpwd == "" || Loginpwd == "undefined") {
			$("#showMeg").html("<i class='icon-err'></i><span style='color:#FF9292;font-size:14px;'>请输入密码 </span>");
			return false;
		} else {
			//   显示加载中
			$(".icon-loading").show();
			//
			var postdata = {
				"Loginuser": Loginuser,
				"Loginpwd": Loginpwd,
				"Os":"PcWEB"
			};
			$.ajax({
				type:"post",
				url:"/login",
				//url:"http://localhost:8080/login",
				async:true,
				data:JSON.stringify(postdata),
				dataType: "JSON",
				success:function(data){
					$(".icon-loading").hide();
					//
					if(data.Rcode == "1000") {
						localStorage.clear();
						localStorage.setItem("LoginUser", JSON.stringify(data.Result))
							//$.widows.location.href="index.html";
						window.location.href = "/web2/html/index.html";
					} else {
						$("#showMeg").html("<i class='icon-err'></i><span style='color:#FF9292;font-size:14px; '>" + data.Reason + " </span>");
					}
				}
			});
		}

		return false;
	});

});
