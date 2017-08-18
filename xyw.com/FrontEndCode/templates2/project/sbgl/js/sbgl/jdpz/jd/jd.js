//开启路由器
app.controller("sbglJdContr", ['$scope', 'httpService', '$modal', '$document', 'alertService', 'toaster', function($scope, httpService, $modal, $document, alertService, toaster) {
	console.log("节点配置-节点");

	//    节点列表
	$scope.jdItems = [];
	$scope.jdxhItems = [];
	//
	$scope.form = {
		//  关键词
		"KeyWord": "",
		//  节点型号ID:字符串  ;   节点型号ID ？？
		"NodeId": "",
		//   教室ID:字符串[可多选]
		"ClassRoomIds": "",
		//   安装位置HTML
		"addHtml": "",
		//   楼层ID:字符串[可多选]
		"Floorsids": "",
		//   楼栋ID:字符串[可多选]
		"Buildingids": "",
		//   校区ID:字符串[可多选]
		"Campusids": "",
		//   是否未安装[0:未选/1:已选]:字符串
		"IsNoSave": "0",
		//   是否未安装bol
		"IsNoSaveBol": false
	}

	//    分页
	$scope.backPage = {
		PageIndex: 1,
		PageSize: 15
	}

	//   设备ID发生变化
	$scope.$watch('form.IsNoSaveBol', function(newValue, oldValue, scope) {
		if(newValue == true) {
			$scope.form.IsNoSave = "1";
		} else {
			$scope.form.IsNoSave = "0";
		}
	});

	//  取节点列表
	var getNodeList = function() {
		var url = config.HttpUrl + "/device/getNodeList";
		var data = {
			"Auth": {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			"Page": {
				"PageIndex": $scope.backPage.PageIndex,
				"PageSize": $scope.backPage.PageSize
			},
			"Para": {
				"KeyWord": $scope.form.KeyWord,
				"NodeId": $scope.form.NodeId,
				"ClassRoomIds": String($scope.form.ClassRoomIds),
				"Floorsids": String($scope.form.Floorsids),
				"Buildingids": String($scope.form.Buildingids),
				"Campusids": String($scope.form.Campusids),
				"IsNoSave": $scope.form.IsNoSave
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取节点列表a", data);
			if(data.Rcode == "1000") {
				$scope.jdItems = data.Result.Data;
				//分页
				data.Result.Page.PageIndex == 0 ? data.Result.Page.PageIndex = 1 : data.Result.Page.PageIndex;
				$scope.backPage = pageFn(data.Result.Page, 5);
			} else {
				console.log(data.Reason);
			}
		});
	}

	//    打开弹窗  选择教室
	$scope.modalOpenClassroom = function() {
		console.log("打开弹窗 -选择教室");
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.itens;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log("qq", selectedItem)
			if(selectedItem.addId == "") {
				$scope.form.Campusids = "";
				$scope.form.Buildingids = "";
				$scope.form.Floorsids = "";
				$scope.form.ClassRoomIds = "";
				$scope.form.addHtml = "";
			} else {
				switch(selectedItem.addCode) {
					case "campus":
						//
						$scope.form.Campusids = selectedItem.addId;
						$scope.form.Buildingids = "";
						$scope.form.Floorsids = "";
						$scope.form.ClassRoomIds = "";
						$scope.form.addHtml = selectedItem.add;
						break;
					case "building":
						//
						$scope.form.Campusids = "";
						$scope.form.Buildingids = selectedItem.addId;
						$scope.form.Floorsids = "";
						$scope.form.ClassRoomIds = "";
						$scope.form.addHtml = selectedItem.add;
						break;
					case "floor":
						//
						$scope.form.Campusids = "";
						$scope.form.Buildingids = "";
						$scope.form.Floorsids = selectedItem.addId;
						$scope.form.ClassRoomIds = "";
						$scope.form.addHtml = selectedItem.add;
						break;
					case "classroom":
						//
						$scope.form.Campusids = "";
						$scope.form.Buildingids = "";
						$scope.form.Floorsids = "";
						$scope.form.ClassRoomIds = selectedItem.addId;
						$scope.form.addHtml = selectedItem.add;
						break;
				}
			}
			$scope.searchPost();
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	};

	/**
	 * 取节点型号列表
	 */
	var getNodeModelList = function() {
		var url = config.HttpUrl + "/device/getNodeModelList";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Page: {
				PageIndex: -1
					//PageSize:PageSize
			},
			Para: {
				//"KeyWord": $scope.form.KeyWord
				"KeyWord": ""
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取节点型号列表", data);
			if(data.Rcode == "1000") {
				$scope.jdxhItems = data.Result.Data;
				$scope.jdxhItems.unshift({
					Id: '',
					Name: '全部节点型号'
				});
			} else {
				console.log(data.Reason);
			}
		});
	}

	/**
	 * 节点删除
	 */
	var deleteNode = function(id) {
		var url = config.HttpUrl + "/device/deleteNode";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				NodeId: id.toString()
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("节点删除", data);
			if(data.Rcode == "1000") {
				toaster.pop('success', '删除成功！');
				getNodeList();
			} else {
				toaster.pop('warning', data.Reason);
			}
		});
	}

	/*  -------------------- 分页、页码  -----------------------  */
	/*----------------
	 //    分页对象添加页码
	 //    return  obj  分页对象
	 //    pagedata:obj  分页对象
	 //    maxpagenumber:int  显示页码数默认5个页码
	 ------------------*/
	var pageFn = function(pagedata, maxpagenumber) {
			if(pagedata.length < 1) return null;
			//   缺省时分5页
			Number(maxpagenumber) > 0 ? maxpagenumber = Number(maxpagenumber) : maxpagenumber = 5;
			var nub = [];
			var mid = Math.ceil(maxpagenumber / 2);
			if(pagedata.PageCount > maxpagenumber) {
				//  起始页
				var Snumber = 1;
				if((pagedata.PageIndex - mid) < 1) {
					Snumber = 1
				} else if((pagedata.PageIndex + mid) > pagedata.PageCount) {
					Snumber = pagedata.PageCount - maxpagenumber + 1;
				} else {
					Snumber = pagedata.PageIndex - (mid - 1)
				}
				for(var i = 0; i < maxpagenumber; i++) {
					nub.push(Snumber + i);
				}
			} else {
				for(var i = 0; i < pagedata.PageCount; i++) {
					nub.push(i + 1);
				}
			}
			pagedata.Number = nub;
			return pagedata;
		}
		//  翻页
	$scope.pageClick = function(pageindex) {
			Number(pageindex) < 1 ? pageindex = 1 : pageindex = Number(pageindex);
			$scope.backPage.PageIndex = pageindex;
			getNodeList();
		}
		/*  -------------------- 分页、页码  -----------------------  */

	//添加按钮功能
	$scope.openModalAddJd = function(str, item) {
			var modalInstance = $modal.open({
				templateUrl: '../project/sbgl/html/sbgl/jdpz/modal_jd.html',
				controller: 'modalSbglJdpzJdContr',
				windowClass: 'm-modal-sbgl-jdpz-jd',
				resolve: {
					items: function() {
						return {
							"operate": str,
							"item": item
						};
					}
				}
			});

			//    弹窗返回
			modalInstance.result.then(function(bol) {
				console.log(bol)
				if(!bol) {
					//
				} else {
					//    刷新列表
					getNodeList();
				}
			}, function() {
				//$log.info('Modal dismissed at: ' + new Date());
			});
		}



	//  删除前提示
	var deleteBefore = function(Id,delFn) {
		var url = config.HttpUrl + "/device/onDeletingNode";
		var data = {
			Auth: {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			Para: {
				NodeId: Id
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("删除前提示", data);
			if(data.Rcode == "1000") {
				delFn(data.Result);
			} else {
				return false;
			}
		});
	}




		//  删除弹窗功能
	$scope.deleteItem = function(item) {
		var modalInstance = $modal.open({
			templateUrl: 'modal/modal_alert_all.html',
			controller: 'modalAlert2Conter',
			resolve: {
				items: function() {
					return {
						"type": 'warning',
						"msg": '确定删除当前节点吗？'
					};
				}
			}
		});
		modalInstance.result.then(function(bul) {
			console.log("bul", bul);
			if(bul) {
				deleteBefore(item.Id, function(rBul) {
					//   是否有关联内容
					if(rBul){
						var modalInstance2 = $modal.open({
							templateUrl: 'modal/modal_alert_all.html',
							controller: 'modalAlert2Conter',
							resolve: {
								items: function() {
									return {
										"type": 'warning',
										"msg": '警告：此操作将影响与当前节点关联的设备控制和联动场景执行！'
									};
								}
							}
						});
						modalInstance2.result.then(function(bul) {
							console.log("bul", bul);
							if(bul) {
								//    删除
								deleteNode(item.Id);
							}
						});
					}else{
						//  删除
						deleteNode(item.Id);
					}
				});
			}
		});

	}




		//  查询
	$scope.searchPost = function() {
			$scope.backPage.PageIndex = 1;
			getNodeList();
			getNodeModelList();
		}

	//   change
	$scope.changeItem = function(item){
		//  select添加时把当前选择的Id赋值
		if(item) {
			$scope.form.NodeId = item.Id
		}
		//   查询
		$scope.searchPost();
	}


		//   回车查询
	$scope.sbgzKeyup = function(e) {
		var keycode = window.event ? e.keyCode : e.which;
		if(keycode == 13) {
			getNodeList();
		}
	}
	$scope.run = function() {
		getNodeList();
		getNodeModelList();
	}
	$scope.run();

}]);

//设备管理-节点配置-节点-添加节点弹窗
app.controller("modalSbglJdpzJdContr", ['$scope', 'httpService', '$modal', '$modalInstance', 'items', "formValidate", 'toaster', function($scope, httpService, $modal, $modalInstance, items, formValidate, toaster) {
	console.log("设备管理-节点配置-节点-添加节点弹窗");
	$scope.items = items;
	console.log("弹窗数据", $scope.items)

	//	取消按钮
	$scope.cancel = function() {
		$modalInstance.dismiss('cancel');
	}

	$scope.form = {
		//    节点编号
		"Id": "",
		//   名称
		"Name": "",
		//   节点型号
		"ModelId": "",
		//   节点型号列表
		"ModelItems": [],
		//   教室
		"ClassRoomId": "",
		//   教室位置
		"ClassRoomHtml": "",
		//   IP类型
		"IpType": "ipv4",
		//   coap
		"NodeCoapPort": "",
		//   路由IP
		"RouteIp": "",
		//   端口
		"InRouteMappingPort": "",
		//   最新上报数据时间
		"UploadTime": ""
	}
	
	//   节点编号失去焦点 验证编号 是否存在
	$scope.blurId = function(){
		
		if($scope.form.Id.length < 1){
			return false;
		}
		
		var url = config.HttpUrl + "/device/getNode";
		var data = {
			"Auth": {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			"Para": {
				//    节点编号
				"NodeId": $scope.form.Id
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("节点编号失去焦点 验证编号 是否存在", data);
			if(data.Rcode == "1000") {
				//   已存在
				formValidate("").minLength(1).outMsg('节点编号已存在，请重新输入！');
			} else {
				//   不存在   通过
				//toaster.pop('warning', data.Reason);
			}
		});
	}
	
	
	

	//    打开弹窗  选择教室
	$scope.modalOpenClassroom = function() {
		console.log("打开弹窗 -选择教室");
		var modalInstance = $modal.open({
			templateUrl: '../html/modal/modal_school.html',
			controller: 'modalGetClassRoomCtrl',
			resolve: {
				items: function() {
					return $scope.itens;
				}
			}
		});

		modalInstance.result.then(function(selectedItem) {
			console.log(selectedItem)
			if(!selectedItem) {
				$scope.form.ClassRoomHtml = "";
				$scope.form.ClassRoomId = "";
			} else {
				if(selectedItem.addCode == "classroom") {
					$scope.form.ClassRoomHtml = selectedItem.add;
					$scope.form.ClassRoomId = selectedItem.addId;
				} else {
					$scope.modalOpenClassroom();
					$modal.open({
						templateUrl: 'modal/modal_alert_all.html',
						controller: 'modalAlert2Conter',
						resolve: {
							items: function() {
								return {
									"type": 'info',
									"msg": '必须选择到教室'
								};
							}
						}
					});
				}
			}
		}, function() {
			//$log.info('Modal dismissed at: ' + new Date());
		});
	};

	//   节点添加
	var saveNode = function() {

		if(!(formValidate($scope.form.Id).minLength(0).outMsg(2407).isOk) || !(formValidate($scope.form.Id).isNode().outMsg(2204).isOk)) return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2408).isOk)) return false;
		if(!(formValidate($scope.form.ClassRoomHtml).minLength(0).outMsg(2409).isOk)) return false;
		var url = config.HttpUrl + "/device/saveNode";
		var data = {
			"Auth": {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			"Para": {
				//    节点编号
				"Id": $scope.form.Id,
				//   名称
				"Name": $scope.form.Name,
				//   节点型号
				"ModelId": $scope.form.ModelId,
				//   教室
				"ClassRoomId": Number($scope.form.ClassRoomId),
				//   IP类型
				"IpType": $scope.form.IpType,
				//   coap
				"NodeCoapPort": $scope.form.NodeCoapPort,
				//   路由IP
				"RouteIp": $scope.form.RouteIp,
				//   端口
				"InRouteMappingPort": $scope.form.InRouteMappingPort,
				//   最新上报数据时间
				"UploadTime": $scope.form.UploadTime
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("节点添加", data);
			if(data.Rcode == "1000") {
				toaster.pop('success', '添加成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		});
	}

	//   节点修改
	var saveNodeEdit = function() {
		if(!(formValidate($scope.form.Id).minLength(0).outMsg(2407).isOk) || !(formValidate($scope.form.Id).isNode().outMsg(2204).isOk)) return false;
		if(!(formValidate($scope.form.ModelId).minLength(0).outMsg(2408).isOk)) return false;
		if(!(formValidate($scope.form.ClassRoomHtml).minLength(0).outMsg(2409).isOk)) return false;
		var url = config.HttpUrl + "/device/saveNode";
		var data = {
			"Auth": {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			"Para": {
				//    节点编号
				"Id": $scope.form.Id,
				//   名称
				"Name": $scope.form.Name,
				//   节点型号
				"ModelId": $scope.form.ModelId,
				//   教室
				"ClassRoomId": Number($scope.form.ClassRoomId),
				//   IP类型
				"IpType": $scope.form.IpType,
				//   coap
				"NodeCoapPort": $scope.form.NodeCoapPort,
				//   路由IP
				"RouteIp": $scope.form.RouteIp,
				//   端口
				"InRouteMappingPort": $scope.form.InRouteMappingPort,
				//   最新上报数据时间
				"UploadTime": $scope.form.UploadTime
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("节点修改", data);
			if(data.Rcode == "1000") {
				toaster.pop('success', '修改成功！');
				$modalInstance.close(true);
			} else {
				toaster.pop('warning', data.Reason);
			}
		});
	}

	//   取节点型号列表
	var getNodeModelList = function() {
		var url = config.HttpUrl + "/device/getNodeModelList";
		var data = {
			"Auth": {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Os": "WEB"
			},
			"Page": {
				"PageIndex": -1
					//"PageSize":$scope.backPage.PageSize
			},
			"Para": {
				"KeyWord": ""
			}
		}
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			console.log("取节点型号列表", data);
			if(data.Rcode == "1000") {
				$scope.form.ModelItems = data.Result.Data;
				for(var item in $scope.form.ModelItems) {
					if($scope.form.ModelItems[item].Id == $scope.form.ModelId) {
						$scope.form.ModelItem = $scope.form.ModelItems[item];
					}
				}
			} else {
        toaster.pop('warning',data.Reason);
			}
		});
	}

	//   生成字符串GUID
	function getGUIDs() {
		var GUID = "";
		for(var i = 1; i <= 32; i++) {
			var n = Math.floor(Math.random() * 16.0).toString(16);
			GUID += n;
			if((i == 8) || (i == 12) || (i == 16) || (i == 20))
			//GUID += "-";
				GUID += "";
		}
		GUID += "";
		return GUID;
	}

	//    打开弹出 -选择日期
	$scope.showDate = function() {
		jeDate({
			dateCell: "#jd_begindate",
			format: "YYYY-MM-DD hh:mm:ss",
			isTime: true,
			minDate: "2015-12-31 00:00:00",
			isinitVal: false,
			choosefun: function(elem, val) {
				$scope.form.UploadTime = val;
			},
			okfun: function(elem, val) {
				$scope.form.UploadTime = val;
			},
			clearfun: function(elem, val) {
				$scope.form.UploadTime = "";
			}
		});
	}

	//  select添加时把当前选择的Id赋值
	$scope.changeModelItem = function(item) {
		$scope.form.ModelId = item.Id;
	}
	$scope.ok = function() {
		//console.log($scope.form)
		//return false;
		//    操作
		switch($scope.items.operate) {
			case "add":
				saveNode();
				break;
			case "look":
				break;
			case "edit":
				saveNodeEdit();
				break;
		}
	}

	$scope.run = function() {
		getNodeModelList();

		switch($scope.items.operate) {
			case "add":
				//$scope.form.Id = getGUIDs();
				break;
			case "look":
				$scope.form.ModelId = $scope.items.item.ModelId;

				$scope.form = $.extend({}, $scope.form, $scope.items.item);
				$scope.form.ClassRoomHtml = $scope.items.item.Campusname + "-" + $scope.items.item.Buildingname + "-" + $scope.items.item.Classroomsname;
				break;
			case "edit":
				$scope.form.ModelId = $scope.items.item.ModelId;

				$scope.form = $.extend({}, $scope.form, $scope.items.item);
				$scope.form.ClassRoomHtml = $scope.items.item.Campusname + "-" + $scope.items.item.Buildingname + "-" + $scope.items.item.Classroomsname;
				break;
		}
	}
	$scope.run();

}]);
