'use strict';

/**
 * 设备管理路由
 * 
 */

app.config(
		['$stateProvider', '$urlRouterProvider',
			function($stateProvider, $urlRouterProvider) {
				$stateProvider
					/*  //////////////////  设备管理   /////////////////////   */
					//   设备管理
						.state('app.sbgl', {
							url: '/sbgl',
							templateUrl: '../project/sbgl/html/sbgl/index.html',
							controller: 'sbglindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../vendor/echarts/echarts.min.js',
													'../project/sbgl/js/sbgl/sbgl.js'
												]
										);
									}
								]
							}
						})
						//  设备型号管理
						.state('app.xhgl', {
							url: '/sbgl/xhgl',
							templateUrl: '../project/sbgl/html/sbgl/sbgl_xhgl.html',
							controller: 'sbglXhglindexContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												['../js/controllers/sbgl.js']
										);
									}
								]
							}
						})
						//   设备日志
						.state('app.sbgl.sbrz', {
							url: '/sbrz',
							templateUrl: '../project/sbgl/html/sbgl/sbrz/index.html',
							controller: 'sbglSbrzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbrz/sbrz.css',
													'../project/sbgl/js/sbgl/sbrz/sbrz.js'
												]
										);
									}
								]
							}
						})


						//   设备故障
						.state('app.sbgl.sbgz', {
							url: '/sbgz',
							templateUrl: '../project/sbgl/html/sbgl/sbgz/index.html',
							controller: 'sbglSbgzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbgz/sbgz.css',
													'../project/sbgl/js/sbgl/sbgz/sbgz.js'
												]
										);
									}
								]
							}
						})

						/*     设备监控         */
						.state('app.sbgl.sbjk', {
							url: '/sbjk',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/index.html',
							controller: 'sbjkContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbjk/sbjk.css',
													'../vendor/jquery/blockUI/jquery.blockUI.js',
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})
						/*     设备监控列表页 教室设备列表         */
						.state('app.sbgl.sbjk.list', {
							url: '/list',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/list.html',
							controller: 'sbjkListContr',
							params:{classroomid:null},
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})

						/*     设备监控列表页 教室设备控制详情         */
						.state('app.sbgl.sbjk.list.article', {
							url: '/article',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/article.html',
							controller: 'sbjkArticleContr',
							params:{DeviceId:null,page:null}
						})

						/*     设备监控列表 查课表         */
						.state('app.sbgl.sbjk.list.ckb', {
							url: '/list',
							templateUrl: '../project/sbgl/html/sbgl/sbjk/schedule.html',
							controller: 'sbjkListCkbContr',
							params:{classroomid:null},
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/js/sbgl/sbjk/sbjk.js'
												]
										);
									}
								]
							}
						})
						/*  /////////////////// 设备监控 End  //////////////////  */
						//   设备预警
						.state('app.sbgl.sbyj', {
							url: '/sbyj',
							templateUrl: '../project/sbgl/html/sbgl/sbyj/index.html',
							controller: 'sbglSbyjContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbyj/sbyj.css',
													'../project/sbgl/js/sbgl/sbyj/sbyj.js'
												]
										);
									}
								]
							}
						})
						//   联动场景
						.state('app.sbgl.ldcj', {
							url: '/ldcj',
							templateUrl: '../project/sbgl/html/sbgl/ldcj/index.html',
							controller: 'sbglLdcjContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/ldcj/ldcj.css',
													'../project/sbgl/js/sbgl/ldcj/ldcj.js'
												]
										);
									}
								]
							}
						})
						//   设备分析
						.state('app.sbgl.sbfx', {
							url: '/sbfx',
							templateUrl: '../project/sbgl/html/sbgl/sbfx/index.html',
							controller: 'sbglSbfxContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbfx/sbfx.css',
													'../project/sbgl/js/sbgl/sbfx/sbfx.js'
												]
										);
									}
								]
							}
						})

						//   节点配置
						.state('app.sbgl.jdpz', {
							url: '/jdpz',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/index.html',
							controller: 'sbglJdpzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jdpz.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-节点型号
						.state('app.sbgl.jdpz.jdxh', {
							url: '/jdxh',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/jdxh/index.html',
							controller: 'sbglJdxhContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jdxh/jdxh.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-控制命令
						.state('app.sbgl.jdpz.kzml', {
							url: '/kzml',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/kzml/index.html',
							controller: 'sbglJdpzKzmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/kzml/kzml.js'
												]
										);
									}
								]
							}
						})

						//   节点配置-节点
						.state('app.sbgl.jdpz.jd', {
							url: '/jd',
							templateUrl: '../project/sbgl/html/sbgl/jdpz/jd/index.html',
							controller: 'sbglJdContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/jdpz/jdpz.css',
													'../project/sbgl/js/sbgl/jdpz/jd/jd.js'
												]
										);
									}
								]
							}
						})


						// --------------  设备配置----------------
						.state('app.sbgl.sbpz', {
							url: '/sbpz',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/index.html',
							controller: 'sbglSbpzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/sbpz.js'
												]
										);
									}
								]
							}
						})

						// --------------  设备配置--故障分类----------------
						.state('app.sbgl.sbpz.gzfl', {
							url: '/gzfl',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/gzfl/index.html',
							controller: 'sbglGzflContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/gzfl/gzfl.js'
													]);
												}
										);
									}]
							}
						})

						// --------------  设备配置--设备型号----------------
						.state('app.sbgl.sbpz.sbxh', {
							url: '/sbxh',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/sbxh/index.html',
//					controller-----控制器
							controller: 'sbglSbxhContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/sbxh/sbxh.js'
													]);
												}
										);
									}]
							}
						})
						// --------------  设备配置--控制命令----------------
						.state('app.sbgl.sbpz.kzml', {
							url: '/kzml',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/kzml/index.html',
//					controller-----控制器
							controller: 'sbglSbpzKzmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/kzml/kzml.js'
												]
										);
									}
								]
							}
						})

						// --------------  设备配置--状态命令----------------
						.state('app.sbgl.sbpz.ztml', {
							url: '/ztml',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/ztml/index.html',
//					controller-----控制器
							controller: 'sbglZtmlContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/ztml/ztml.js'
												]
										);
									}
								]
							}
						})
						// --------------  设备配置--故障现象----------------
						.state('app.sbgl.sbpz.gzxx', {
							url: '/gzxx',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/gzxx/index.html',
//					controller-----控制器
							controller: 'sbglGzxxContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function( $ocLazyLoad ){
										return $ocLazyLoad.load('angularBootstrapNavTree').then(
												function(){
													return $ocLazyLoad.load([
														'../project/sbgl/css/sbgl/sbpz/sbpz.css',
														'../project/sbgl/js/sbgl/sbpz/gzxx/gzxx.js'
													]);
												}
										);
									}]
							}
						})
						/*   --------设备配置--状态命令-- 添加 code     */
						.state('app.sbgl.sbpz.ztml.code', {
							url: '/code',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/ztml/code/index.html',
							controller: 'modalZtmlAddContr',
							params: {
								mid:"",
								code: "",
								code_name:""
							},
							resolve: {
								deps: ['$ocLazyLoad',
									function ($ocLazyLoad) {
										return $ocLazyLoad.load(['ui.select']).then(
												function () {
													return $ocLazyLoad.load(
															[
																'../project/sbgl/css/sbgl/sbpz/sbpz.css',
																'../project/sbgl/js/sbgl/sbpz/ztml/ztml.js'
															]
													);
												}
										);
									}
								]
							}
						})
						// --------------  设备配置--设备管理者----------------
						.state('app.sbgl.sbpz.sb', {
							url: '/sbglz',
							templateUrl: '../project/sbgl/html/sbgl/sbpz/sbglz/index.html',
//					controller-----控制器
							controller: 'sbglSbglzContr',
							resolve: {
								deps: ['$ocLazyLoad',
									function($ocLazyLoad) {
										return $ocLazyLoad.load(
												[
													'../project/sbgl/css/sbgl/sbpz/sbpz.css',
													'../project/sbgl/js/sbgl/sbpz/sbglz/sbglz.js'
												]
										);
									}
								]
							}
						})
						/*  //////////////////  设备管理 End /////////////////////   */







			}
		]
	);