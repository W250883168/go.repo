/**
 * ---------------------------------------------
 *
 * (c) Copyright 2017 wu. All Rights Reserved.
 *
 * 后台管理系统指令
 *
 * ---------------------------------------------
 */


/*    弹窗指令     */
app.directive('alertBar', [function() {
	return {
		restrict: 'EA',
		templateUrl: '../project/zkmb/html/zkmb/modal_alert.html',
		scope: {
			message: "=",
			type: "="
		},
		link: function(scope, element, attrs) {

			scope.hideAlert = function() {
				scope.message = null;
				scope.type = null;
			};

		}
	};
}]);





//    百度 echarts 指令
app.directive('eChart', [function () {

    function link($scope, element, attrs) {
        // 基于准备好的dom，初始化echarts图表
        var myChart = echarts.init(element[0]);

        //监听options变化
        if (attrs.uiOptions) {
            attrs.$observe('uiOptions', function () {
                var options = $scope.$eval(attrs.uiOptions);
                if (angular.isObject(options)) {
                	if(options.fn){
                		//  formatter
	                	options = $scope.$parent.formatter(options,element);
	                };
                    myChart.setOption(options);
                }
            }, true);
        }
    }
    return {
        restrict: 'A',
        link: link
    };
}]);


//    图片上传 指令
app.directive('upImg', ['httpService','toaster','$modal',function (httpService,toaster,$modal) {

    function link($scope, element, attrs) {
    	//element.
    	element.find("input[type=file]").on('change',function(){
    		//   验证上传图片大小  2M或内
    		var size = $(this)[0].files[0].size;
    		if(size > 1024 * 1024 * 2){
    			$modal.open({
	            	templateUrl: 'modal/modal_alert_all.html',
	            	controller: 'modalAlert2Conter',
	            	resolve: {
	            		items: function () {
	            			return {"type":'danger',"msg":'图片大小必须2M或已内！'};
	            		}
	            	}
	            });
	            $(this)[0].value = '';
	            return;
    		}
    		
    		var action = config.HttpUrl + '/upfile';
    		//    有自定义dom ID  ，  form action 。
    		if("action" in attrs && "id" in attrs){
    			action = config.HttpUrl + '/vod/attachment/upload';
    			//   input 其他属性
    			if("input" in attrs){
    				var dataInput = $scope.$eval(attrs.input);
    				var htmlInput = "";
    				for(var a in dataInput){
    					htmlInput += '<input type="text" name="' + a + '" value="' + dataInput[a] + '" />';
    				}
    			}
    		}

    		//生成空白的iframe
		    var iframe = $('<iframe src="javascript:false;" />');
		    //掩藏iframe
		    iframe.hide();
		    element.append(iframe);
		    var temp = element.find("input[type=file]").clone(true,true);
		    var iframe_body = $(iframe[0].contentWindow.document.body);
		    iframe_body.append('<form method="post" enctype="multipart/form-data" action="' + action + '"></form>');
		    iframe_body.find("form").append(element.find("input[type=file]"));
		    iframe_body.find("form").append($("<input name='filesize' value='" + size + "' type='hidden' />"));
		    //    有自定义dom ID  ，  form action 。
    		if("action" in attrs && "id" in attrs){
    			if("input" in attrs){
    				//   form 添加INPUT
					iframe_body.find("form").append(htmlInput);
    			}
    		}
    		//element.append(temp);
    		element.append(temp).append("<i class='icon-refresh'></i>").find(">span").text("上传中");
    		iframe_body.find("form").submit();
    		//
    		iframe.load(function(e){
    			//
    			element.find(">span").text('上传');
    			element.find("i.icon-refresh").remove();
    			//
    			var body = $(e.target.contentWindow.document.body);
    			var img_item = body.find("pre").text();
    			if(img_item != ""){
    				var temp = $scope.$eval(img_item);
    				//element.hide();
    				 $scope.$apply(function () {
    				 	if("action" in attrs && "id" in attrs){
    				 		$scope[attrs.id] ? $scope[attrs.id].push(temp) : $scope[attrs.id] = [temp];
    				 	}else{
    				 		$scope.upimglist.push(temp);
    				 	}
                    });
    			}
    		});
    	});


    	//   图片删除
    	$scope.removeFile = function(item,index){
    		/**
			 * 删除图片
			 */
			var deleleattachment = function(id) {
				var url = config.HttpUrl + "/vod/deleleattachment";
				var data = {
		            "Usersid": config.GetUser().Usersid,
		            "Rolestype": config.GetUser().Rolestype,
		            "Token": config.GetUser().Token,
		            "Os": "WEB",
		            "ID":Number(id)
				};
				var promise = httpService.ajaxPost(url, data);
				promise.then(function(data) {
					console.log("删除图片",data)
					if(data.Rcode == "1000") {
						$scope[attrs.id].splice(index, 1);
					} else {
            toaster.pop('warning',data.Reason);
					}
				}, function(reason) {}, function(update) {});
			}
			//  run
			deleleattachment(item.ID);
    	}

    }
    return {
        restrict: 'A',
        link: link
    };
}]);



//    面包屑指令
app.directive("wBreadcrumb",['$rootScope',function($rootScope){

	function link($scope,element,attrs){
		//   不在模板中的路由面包屑
		$scope.breadcrumbAddTitle = attrs.operatetitle;

		$scope.$on('$stateChangeSuccess', function (ev, to, toParams, from, fromParams) {
			//   监控 navlist 取完成
			$scope.$watch('navlist',function(newValue,oldValue, scope){
				$scope.breadcrumb = bre(to.name,newValue,$scope);
			});
    	});
	}

	//var tempHtml =


	//   面包屑组合
	//   str:app.sbgl.sbfx
	//   items:obj   getapp  返回数组数据
	function bre(str,items,$scope){
		if(!str)return false;
		//   app.sbgl.sbfx 等 转数组 容器
		var arr = [];
		arr = str.split(".");

		if(arr.length > 1){
			//   连接后字符串 变量
			var temp = "";

			//   面包屑 容器
			var breadcrumb = [];
			for(var i = 1; i < arr.length; i++){
				temp == "" ? temp += arr[i] : temp += "." + arr[i];
				var isSon = false;
				var tempName = "";
				//
				for(var a in items){
					if(items[a].Modulecode == temp){
						tempName = items[a].Modulename;
						//    是否是干节点模块
						isSon = $scope.haveSon(items[a].Id);
					}
				}
				//   找到有
				if(tempName){
					//  最后一个、当前模块不用链接hide==true;//   干节点不用链接
					if(arr.length - 1 == i || isSon){
						breadcrumb.push({"Modulecode":"app." + temp,"Modulename":tempName,'hide':true});
					}else{
						breadcrumb.push({"Modulecode":"app." + temp,"Modulename":tempName});
					}
				}
			}
		}
		return breadcrumb;
	}

	return {
		restrict: 'EAC',
		link:link,
		template:"<li ng-repeat='item in breadcrumb' ng-class={'active':item.hide}>" +
				 	"<a ng-if='item.hide != true' ui-sref='{{item.Modulecode}}'>{{item.Modulename}}</a>" +
				 	"<span ng-if='item.hide'>{{item.Modulename}}</span>" +
				 "</li>" +
				 "<li ng-if='breadcrumbAddTitle'><span>{{breadcrumbAddTitle}}</span></li>"
	}
}]);






//    提交按键点击
app.directive('submitTimeout', ['$timeout',function ($timeout) {

    function link($scope, element, attrs) {
    	element.on("click",function(){
    		element.attr("disabled","disabled");
    		$timeout(function() {
				element.removeAttr('disabled');
			}, 2000);
    	});
    }
    return {
        restrict: 'A',
        link: link
    };
}]);


//    视频
app.directive('uiVideo', ['$timeout',function ($timeout) {

    function link($scope, element, attrs) {
    	if($scope.items.videoUrl){
    		var videoHtml = '' +
	    		'<video class="splb_video" width="100%" height="auto" controls="controls">' +
					'<source src="' + $scope.items.videoUrl + '" type="video/mp4" />' +
				'</video>';
	    	element.append(videoHtml);
    	}
    }
    return {
        restrict: 'A',
        link: link
    };
}]);




//
app.directive('modalMove', ['$timeout','$document', function($timeout,$document) {
    function link(scope, element, attr) {
    	//   延迟生成dom
    	var timer = $timeout(function(){
    		//
    		var dialog = element.parents(".modal-dialog");
    		var dialogHeight = dialog.height();
    		var wHeight = jQuery(window).height();
    		var wWidth = jQuery(window).width();
    		var header = dialog.find(".g-modal-header");
    		//   弹窗居中
    		if(dialogHeight < wHeight){
    			dialog.css("margin-top",(wHeight - dialogHeight) / 2);
    		}
    		//
    		var startX = 0, startY = 0, x = 0, y = 0;
            //
            header.css({
                position: 'relative',
                cursor: 'move'
            });
			//
            header.on('mousedown', function(event) {
                // Prevent default dragging of selected content
                event.preventDefault();
                startX = event.pageX - x;
                startY = event.pageY - y;
                $document.on('mousemove', mousemove);
                $document.on('mouseup', mouseup);
            });

            function mousemove(event) {
                y = event.pageY - startY;
                x = event.pageX - startX;
                //   不超过屏幕
                if(event.pageY > 0 && event.pageX > 0 && event.pageX < wWidth && event.pageY < wHeight){
                	dialog.css({
	                top: y + 'px',
	                left:  x + 'px'
	                });
                }
            }

            function mouseup() {
                $document.off('mousemove', mousemove);
                $document.off('mouseup', mouseup);
            }


    	},50);
    };
    return {
        restrict: 'A',
        link: link
    };
}]);










//    教室导流 - 导流分析 - 百度 echarts 指令
app.directive('topMap', [function () {

    function link($scope, element, attrs) {
        // 基于准备好的dom，初始化echarts图表
        var myChart = echarts.init(element[0]);

        //监听options变化
        if (attrs.uiScopename) {
            $scope.$watch(attrs.uiScopename + ".chartOptionBul", function () {
                if (angular.isObject($scope[attrs.uiScopename])) {
                    myChart.setOption($scope[attrs.uiScopename]);
                    console.log($scope[attrs.uiScopename])
                }
            }, true);
        }
    }
    return {
        restrict: 'A',
        link: link
    };
}]);




//    分页指令
app.directive("fooPagination",['$rootScope',function($rootScope){
	//   分页对象变量名
	var b_n = "backPage";
	function link($scope,element,attrs){
		//   不在模板中的路由面包屑
		if(attrs.scope){
			b_n = attrs.scope;
		}
	}

	return {
		restrict: 'EAC',
		link:link,
		template:'<ul class="pagination2 clearfix" ng-if="' + b_n + '.PageCount > 1">'
				+	'<li class="footable-page-arrow" ng-class="{disabled:' + b_n + '.PageIndex <= 1}">'
				+		'<a ng-click="pageClick(' + b_n + '.PageIndex -1)">‹</a>'
				+	'</li>'
				+	'<li class="footable-page" ng-if="' + b_n + '.Number[0] > 1">'
				+		'<a ng-click="pageClick(1)">1</a>'
				+	'</li>'
				+	'<li class="footable-page" ng-if="' + b_n + '.Number[0] > 2  && ' + b_n + '.Number.length < ' + b_n + '.PageCount">'
				+		'<span>...</span>'
				+	'</li>'
				+	'<li class="footable-page" ng-repeat="item in ' + b_n + '.Number" ng-class="{active:' + b_n + '.PageIndex == item}">'
				+		'<a  ng-click="pageClick(item)">{{item}}</a>'
				+	'</li>'
				+	'<li class="footable-page" ng-if="' + b_n + '.Number[' + b_n + '.Number.length - 1] < ' + b_n + '.PageCount - 1">'
				+		'<span>...</span>'
				+	'</li>'
				+	'<li class="footable-page" ng-if="' + b_n + '.Number[' + b_n + '.Number.length - 1] != ' + b_n + '.PageCount">'
				+		'<a ng-click="pageClick(' + b_n + '.PageCount)">{{' + b_n + '.PageCount}}</a>'
				+	'</li>'
				+	'<li class="footable-page-arrow" ng-class="{disabled:' + b_n + '.PageIndex >= ' + b_n + '.PageCount}">'
				+		'<a ng-click="pageClick(' + b_n + '.PageIndex + 1)">›</a>'
				+	'</li>'
				+'</ul>'
				+'&nbsp;&nbsp;&nbsp;'
				+'<ul class="pagination2 clearfix" ng-if="' + b_n + '.PageCount > 1 && ' + b_n + '.Number.length < ' + b_n + '.PageCount">'
				+	'<li class="footable-page" ng-init="temp_backPage_index = 1"><input type="number" ng-model="temp_backPage_index" /></li>'
				+	'<li class="footable-page">'
				+		'<a ng-click="pageClick( ( (temp_backPage_index <= ' + b_n + '.PageCount) && (temp_backPage_index >= 1) ) ? temp_backPage_index : ' + b_n + '.PageCount)">跳转</a>'
				+	'</li>'
				+'</ul>'
	}
}]);
