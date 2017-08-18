'use strict';
/**
 * Created by Administrator on 2016/7/28.
 * 用户权限分配
 */

/*  //////////////////////////////////////////////////  */
//   用户权限分配
app.controller("qxglQxfpContr", ['$scope', 'httpService', '$filter', '$state', 'toaster', '$modal', function($scope, httpService, $filter, $state, toaster, $modal) {
	//    select
	$scope.form = {
		"jsgl_listitem": "",
		//  ----
		"jsgl_listitem_item": ""
	}

	//   模块树
	$scope.mkgl_tree = {};
	//  权限分配树-模块树
	$scope.mkgl_tree_qxfp = [];
	//  模块数据表转树
	function convert(source) {
		var tmp = {},
			parent, n;
		for(n in source) {
			var item = source[n];
			if(item.Id == item.Superiormoduleid) {
				parent = item.Id;
			}
			if(!tmp[item.Id]) {
				tmp[item.Id] = {};
			}

			tmp[item.Id].Id = item.Id;
			tmp[item.Id].Modulename = item.Modulename;
			tmp[item.Id].Modulecode = item.Modulecode;
			tmp[item.Id].Moduleicon = item.Moduleicon;
			tmp[item.Id].Moduleurl = item.Moduleurl;
			tmp[item.Id].Moduleattribute = item.Moduleattribute;
			tmp[item.Id].Superiormoduleid = item.Superiormoduleid;
			tmp[item.Id].Functionlist = [];

			if(!("children" in tmp[item.Id])) tmp[item.Id].children = [];

			if(item.Id != item.Superiormoduleid) {
				if(tmp[item.Superiormoduleid]) {
					tmp[item.Superiormoduleid].children.push(tmp[item.Id]);
				} else {
					tmp[item.Superiormoduleid] = {
						children: [tmp[item.Id]]
					};
				}
			}
		}

		tmp[0].Id = 0;
		tmp[0].checkbox = false;
		tmp[0].ban = false;
		tmp[0].Modulename = '根目录';

		return tmp;

	}

	//  角色
	//  角色 列表数组
	$scope.jsgl_list = [];

	var init_data = function() {
		var url = config.HttpUrl + "/system/sm/getroles";

		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Rolesname": "",
			"Id": 0,
			"PageIndex": 1,
			"PageSize": 100
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.jsgl_list = data.Result.PageData;
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	//  功能 中间权限查询
	$scope.zjqx_list = {};
	//  选择角色
	$scope.qxgleselect = function(id) {
		//   选中
		$scope.form.jsgl_listitem = id;
		
		//   选择空
		if(!id) {
			//  清空模块、功能选中项
			//			gngl_ck_null();
			//			mkgl_ck_null();
			clear_tree_check($scope.mkgl_tree_qxfp[0]);
			return;
		}
		//  截入模块 。
		var url = config.HttpUrl + "/system/sm/getsetsystemconfigall";
		var Init_id = Number(id);
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Id": Init_id
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			//console.log("选择角色",data)
			if(data.Rcode == "1000") {
				$scope.zjqx_list = data.Result;

				//  清空模块、功能选中项
				//				gngl_ck_null();
				//				mkgl_ck_null();
				clear_tree_check($scope.mkgl_tree_qxfp[0]);
				//  载入初始标记
				//rei_null($scope.zjqx_list);
				Load_dataval($scope.zjqx_list, $scope.mkgl_tree_qxfp[0]);
				//  半选
				ban_bul();
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	}

	var rmfclist = []; //需要添加的角色模块功能中间数据
	var rmclist = []; //需要添加的角色模块中间数据
	var delrmfclist = []; //删除的角色模块功能中间表数据
	var delrmclist = []; //删除的角色模块中间表数据
	//获取需要删除的角色模块中间表数据等数据
	var GetDelValData = function(m, g, tree) {
		for(var index = 0; index < tree.children.length; index++) {
			if(tree.children[index].children.length > 0) { //判断是否有下级
				GetDelValData(m, g, tree.children[index]);
			}
			for(var mindex = 0; mindex < m.length; mindex++) {
				if(tree.children[index].Id == m[mindex].Systemmoduleid) { //判断
					if(!tree.children[index].checkbox) {
						delrmclist.push({
							Id: m[mindex].Id
						});
					}
				}
			}
			for(var gindex = 0; gindex < g.length; gindex++) {
				for(var te = 0; te < tree.children[index].Functionlist.length; te++) {
					if(g[gindex].Systemmodulefunctionsid == tree.children[index].Functionlist[te].Id) { //判断选择的功能项
						if(!tree.children[index].Functionlist[te].checkbox) { //判断是否选中[未选择则将删除的数据加入进去]
							delrmfclist.push({
								Id: g[gindex].Id
							}); //.push(g[gindex].Id);
						}
					}
				}
			}
		}
	};
	//获取新加入设置的值
	var GetSetValData = function(tree, rid) {
		for(var index = 0; index < tree.children.length; index++) {
			if(tree.children[index].children.length > 0) { //判断是否有下级
				GetSetValData(tree.children[index], rid);
			}
			console.log(tree.children[index].checkbox);
			if(tree.children[index].checkbox) {
				var state = 0;
				if(tree.children[index].ban) {
					state = 1;
				}
				var isadd = true;
				for(var rmcindex = 0; rmcindex < rmclist.length; rmcindex++) {
					if(rmclist[rmcindex].Systemmoduleid == tree.children[index].Id) {
						isadd = false;
						break;
					}
				}
				if(rmclist.length == 0 || isadd) {
					rmclist.push({
						Rolesid: rid,
						Systemmoduleid: tree.children[index].Id,
						State: state
					}); //角色模块中间数据
				}
			}
			for(var te = 0; te < tree.children[index].Functionlist.length; te++) {
				if(tree.children[index].Functionlist[te].checkbox) { //判断是否选中[未选择则将删除的数据加入进去]
					//					rmfclist.push({Systemmodulefunctionsid:tree.children[index].Functionlist[te].Id});
					var isadd = true;
					for(var rmfcindex = 0; rmfcindex < rmfclist.length; rmfcindex++) {
						if(rmfclist[rmfcindex].Systemmoduleid == tree.children[index].Functionlist[te].Id) {
							isadd = false;
							break;
						}
					}
					if(rmclist.length == 0 || isadd) {
						rmfclist.push({
							Systemmodulefunctionsid: tree.children[index].Functionlist[te].Id
						}); //角色模块中间数据
					}
				}
			}
		}
	};

	//   提交与修改
	$scope.qxfp_post = function() {
		if(!$scope.form.jsgl_listitem) {
			$modal.open({
				templateUrl: 'modal/modal_alert_all.html',
				controller: 'modalAlert2Conter',
				resolve: {
					items: function() {
						return {
							"type": 'info',
							"msg": '请选择角色？'
						};
					}
				}
			});
			return false;
		}
		rmfclist = []; //需要添加的角色模块功能中间数据
		rmclist = []; //需要添加的角色模块中间数据
		delrmfclist = []; //删除的角色模块功能中间表数据
		delrmclist = []; //删除的角色模块中间表数据
		GetDelValData($scope.zjqx_list[0], $scope.zjqx_list[1], $scope.mkgl_tree_qxfp[0]);
		GetSetValData($scope.mkgl_tree_qxfp[0], Number($scope.form.jsgl_listitem));
		/*
		//		//  模块ID功能对应模块ID
		//		var rmclist_mk_id = [];
		//		//   rmclist 模块中间
		//		dg_tree($scope.mkgl_tree_qxfp, function(item) {
		//			if(item.Id != 0 && !(!(item.oricheckbox)) && item.checkbox == false) {
		//				delrmclist.push({"Id": Number(item.rmclistid)});
		//			} else if(item.Id != 0 && !(item.oricheckbox) && item.checkbox == true) {
		//				rmclist.push({
		//					"Rolesid": Number($scope.form.jsgl_listitem),
		//					"Systemmoduleid": Number(item.Id)
		//				});
		//			}
		//		});
		//		//   rmfclist  功能中间
		//		for(var i = 0; i < $scope.gngl_list.length; i++) {
		//			if(!(!($scope.gngl_list[i].oricheckbox)) && $scope.gngl_list[i].checkbox == false) {
		//				//   原来选中，现在未选中 ，为删除
		//				delrmfclist.push({
		//					"Id": Number($scope.gngl_list[i].rmfclistid)
		//				});
		//			} else if(!($scope.gngl_list[i].oricheckbox) && $scope.gngl_list[i].checkbox == true) {
		//				//  原来没选中，现在选中，添加
		//				rmfclist.push({
		//					"Systemmodulefunctionsid": Number($scope.gngl_list[i].Id)
		//				});
		//				rmclist_mk_id.push({
		//					"Rolesid": Number($scope.form.jsgl_listitem),
		//					"Systemmoduleid": Number($scope.gngl_list[i].Systemmoduleid)
		//				});
		//			}
		//		}
		//		//   单功能添加处理
		//		if(rmfclist.length > 0 && rmclist.length == 0) {
		//			rmclist = rmclist_mk_id;
		//		}
		*/
		//return false;
		var url = config.HttpUrl + "/system/sm/setsystemconfig";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Rmfclist": rmfclist, //需要添加的角色模块功能中间数据
			"Rmclist": rmclist, //需要添加的角色模块中间数据
			"DelRmfclist": delrmfclist, //删除的角色模块功能中间表数据
			"DelRmclist": delrmclist //删除的角色模块中间表数据
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				//  重新载入select值的树
				toaster.pop('success', '提交成功！');
				$scope.qxgleselect($scope.form.jsgl_listitem);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};

	/*  /////////////////////////////////////////  */

	//  清空模块
	var mkgl_ck_null = function() {
			dg_tree($scope.mkgl_tree_qxfp, function(item) {
				item.checkbox = false;
				item.rmclistid = ""; //角色模块中间表Id
				item.oricheckbox = false;
			});
		}
		//   清空功能
	var gngl_ck_null = function() {
			for(var b = 0; b < $scope.gngl_list.length; b++) {
				$scope.gngl_list[b].checkbox = false;
				$scope.gngl_list[b].rmfclistid = ""; //角色模块功能中间表ID
				$scope.gngl_list[b].oricheckbox = false;
			}
		}
		//$scope.mkgl_tree_qxfp
		//清除选择的事项
	var clear_tree_check = function(list) {
		for(var index = 0; index < list.children.length; index++) {
			if(list.children[index].children.length > 0) { //判断是否有下级
				clear_tree_check(list.children[index]);
			}
			list.children[index].checkbox = false;
			for(var findex = 0; findex < list.children[index].Functionlist.length; findex++) { //循环清除功能项的选择
				list.children[index].Functionlist[findex].checkbox = false;
			}
		}
	};
	//加载选择数据
	var Load_dataval = function(data, tree) {
		var mlist = data[0];
		var glist = data[1];
		list_loaddata(mlist, glist, tree);
	};
	//树形分配选择数据
	var list_loaddata = function(m, g, tree) {
		for(var index = 0; index < tree.children.length; index++) {
			if(tree.children[index].children.length > 0) { //判断是否有下级
				list_loaddata(m, g, tree.children[index]);
			}
			for(var mindex = 0; mindex < m.length; mindex++) {
				if(tree.children[index].Id == m[mindex].Systemmoduleid) {
					//$scope.click_checkbox(tree.children[index]);
					if(m[mindex].State == 0) {
						tree.children[index].checkbox = true;
					} else {
						tree.children[index].ban = true;
					}
				}
			}
			for(var gindex = 0; gindex < g.length; gindex++) {
				for(var te = 0; te < tree.children[index].Functionlist.length; te++) {
					if(g[gindex].Systemmodulefunctionsid == tree.children[index].Functionlist[te].Id) {
						tree.children[index].Functionlist[te].checkbox = true;
					}
					//					console.log(tree.children[index].Functionlist[te]);
				}
			}
		}
	};

	//  模块默认是否展开   true 收起， false 展开
	var rei_item_bul = function(items) {
		//  递归设置
		dg_tree(items, function(item) {
			//  设置全部收起
			item.bul = true;
		});
	}

	//  初始化原始标记
	var rei_null = function(zjqx) {

		//　　模块
		for(var i = 0; i < zjqx[0].length; i++) {
			//  递归设置
			dg_tree($scope.mkgl_tree_qxfp, function(item) {
				if(item.Id == zjqx[0][i].Systemmoduleid) {
					item.rmclistid = zjqx[0][i].Id;
					item.checkbox = true;
					//  标记原始选中否
					item.oricheckbox = true;
					return;
				} else {
					//item.checkbox = false;
					//  标记原始选中否
					//item.oricheckbox = false;
					//return;
				}
			});
		}

		//  功能
		for(var i = 0; i < zjqx[1].length; i++) {
			for(var b = 0; b < $scope.gngl_list.length; b++) {
				if(zjqx[1][i].Systemmodulefunctionsid == $scope.gngl_list[b].Id) {
					$scope.gngl_list[b].rmfclistid = zjqx[1][i].Id;
					$scope.gngl_list[b].checkbox = true;
					//  标记原始选中否
					$scope.gngl_list[b].oricheckbox = true;
				} else {
					//$scope.gngl_list[b].checkbox = false;
					//  标记原始选中否
					//$scope.gngl_list[b].oricheckbox = false;
				}
			}
		}
	}

	//  递归树
	var dg_tree = function(tree, fn) {
		for(var i = 0; i < tree.length; i++) {
			fn(tree[i]);
			if(tree[i].children.length > 0) {
				dg_tree(tree[i].children, fn);
			}
		}
	}

	//  去重
	var cast_a = function(at) {
			var res = [];
			var json = {};
			for(var i = 0; i < at.length; i++) {
				if(!json[at[i]]) {
					res.push(at[i]);
					json[at[i]] = 1;
				}
			}
			return res;
		}
		/*  /////////////////////////////////////////  */

	//  功能选择
	$scope.click_gngl = function(item) {
		dg_tree($scope.mkgl_tree_qxfp, function(ite) {
			if(ite.Id == item.Systemmoduleid) {
				if(item.Functionsattribute == "1") {
					for(var i = 0; i < $scope.gngl_list.length; i++) {
						if($scope.gngl_list[i].Systemmoduleid == item.Systemmoduleid) {
							item.checkbox = false;
							ite.checkbox = false;
							$scope.gngl_list[i].checkbox = false;
						}
					}
				}
				//  选中树     子节点全选 全不选
				if(ite.children.length > 0) {
					//  遍历 子节点
					dg_tree([ite], function(ite2) {
						if(ite.checkbox) {
							ite2.checkbox = true;
						} else {
							ite2.checkbox = false;
						}
					});
				}
			}
		});
		//$scope.click_checkbox(item);
	}

	//$scope

	//  选择
	$scope.click_checkbox = function(item) {
		//  选中树     子节点全选 全不选
		if(item.children.length > 0) {
			//  遍历 子节点
			dg_tree([item], function(ite) {
				if(item.checkbox) {
					ite.checkbox = true;
				} else {
					ite.checkbox = false;
				}
				if(ite.Functionlist != undefined) {
					for(var findex = 0; findex < ite.Functionlist.length; findex++) {
						ite.Functionlist[findex].checkbox = ite.checkbox;
					}
				}
			});
		}
		if(item.Functionlist != undefined) {
			//console.log(item.Functionlist.length);
			for(var findex = 0; findex < item.Functionlist.length; findex++) {
				item.Functionlist[findex].checkbox = item.checkbox;
				console.log(item.Functionlist[findex].checkbox);
			}
		}
		//  半选
		ban_bul();
	};

	//  半选   遍历半选 显示
	var ban_bul = function() {
			//  半选
			dg_tree($scope.mkgl_tree_qxfp, function(ite) {
				//   bul_ban_c = false && bul_ban = true 为 全没选。  两个都为true 为半先。
				var bul_ban = false;
				var bul_ban_c = false;
				var bul_ban_b = false;
				var ban = function(tree_i) {
					for(var b = 0; b < tree_i.children.length; b++) {
						if(tree_i.children[b].checkbox) {
							bul_ban_c = true;
						} else if(tree_i.children[b].ban) {
							bul_ban_b = true;
						} else {
							bul_ban = true;
						}
						if(tree_i.children[b].length > 0) {
							ban(tree_i.children[b].children);
						}
					}
				}
				ban(ite);
				if((bul_ban || bul_ban_b) && bul_ban_c) {
					//  半选中
					ite.ban = true;
				} else if(bul_ban_c && bul_ban == false) {
					//  全选
					ite.checkbox = true;
					ite.ban = false;
				} else if(bul_ban_c == false && bul_ban == true) {
					//  全不选
					ite.checkbox = false;
					ite.ban = false;
				}
			});
		}
		//  选择
	$scope.fclick_checkbox = function(item, mitem) {
		//		//  选中树     子节点全选 全不选
		//		if(item.children.length > 0) {
		//			//  遍历 子节点
		//			dg_tree([item], function(ite) {
		//				if(item.checkbox) {
		//					ite.checkbox = true;
		//				} else {
		//					ite.checkbox = false;
		//				}
		//			});
		//		}
		if(item.checkbox) {
			item.checkbox = true;
			mitem.checkbox = true;
		} else {
			item.checkbox = false;
			mitem.checkbox = false;
			for(var mindex = 0; mindex < mitem.Functionlist.length; mindex++) {
				if(mitem.Functionlist[mindex].checkbox) {
					mitem.checkbox = true;
				}
			}
		}
		//  半选
		ban_bul();
	}

	//  是否叶子节点
	$scope.isLeaf = function(id) {
		var bol = true;
		for(var i = 0; i < $scope.mkgl_list.length; i++) {
			if($scope.mkgl_list[i].Superiormoduleid == id) {
				bol = false;
				break;
			}
		}
		return bol;
	}

	//  点击active
	$scope.isOpen = function(item) {
		item.bul = !item.bul;
	}

	//  模块 列表数组
	$scope.mkgl_list = [];

	var init_data_mkgl = function() {
			var url = config.HttpUrl + "/system/sm/getsystemmodel";
			var data = {
				"Usersid": config.GetUser().Usersid,
				"Rolestype": config.GetUser().Rolestype,
				"Token": config.GetUser().Token,
				"Modulename": "",
				"Superiormoduleid": -1,
				"Id": 0,
				"PageIndex": 1,
				"PageSize": 100
			};
			var promise = httpService.ajaxPost(url, data);
			promise.then(function(data) {
				if(data.Rcode == "1000") {
					$scope.mkgl_list = data.Result.PageData;
					//console.log($scope.mkgl_list)
					//  数据转树
					$scope.mkgl_tree = convert($scope.mkgl_list);
					//  权限分配HTML 列表用树数组
					$scope.mkgl_tree_qxfp = [$scope.mkgl_tree[0]];
					//  设置树默认收起
					rei_item_bul($scope.mkgl_tree_qxfp[0].children);
					//console.log("tree",$scope.mkgl_tree_qxfp);
					init_data_gngl();
				} else {
					toaster.pop('warning', data.Reason);
				}
			}, function(reason) {}, function(update) {});
		}
	

	//  功能 列表数组
	$scope.gngl_list = [];

	var init_data_gngl = function() {
		var url = config.HttpUrl + "/system/sm/getsystemmodelfunc";
		var data = {
			"Usersid": config.GetUser().Usersid,
			"Rolestype": config.GetUser().Rolestype,
			"Token": config.GetUser().Token,
			"Functionname": "",
			"Id": 0,
			"PageIndex": 1,
			"PageSize": 10000
		};
		var promise = httpService.ajaxPost(url, data);
		promise.then(function(data) {
			if(data.Rcode == "1000") {
				$scope.gngl_list = data.Result.PageData;
				Mergeprocesslist($scope.mkgl_tree_qxfp[0],$scope.gngl_list);
				//console.log($scope.gngl_list);
			} else {
				toaster.pop('warning', data.Reason);
			}
		}, function(reason) {}, function(update) {});
	};
	
	
	//   放入功能
	var Mergeprocesslist = function(m, g) {
		//   遍历添加功能
		dg_tree(m.children, function(ite) {
			for(var gi = 0; gi < g.length; gi++) {
				//   限制到叶子节点才压入功能
				if(ite.children.length == 0 && ite.Id == g[gi].Systemmoduleid){
					var isadd = true;
					for(var flt = 0; flt < ite.Functionlist.length; flt++) {
						if(ite.Functionlist[flt].Id == g[gi].Id){
							isadd = false;
							break;
						}
					}
					if(isadd) {
						ite.Functionlist.push(g[gi]);
					}
				}
			}
		});
		
	};
	
	
	

	//  run
	var run = function() {
		//  取角色 
		init_data();
		//  run
		init_data_mkgl();

	}
	run();
}]);