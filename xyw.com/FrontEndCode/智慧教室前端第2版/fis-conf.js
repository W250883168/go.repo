// default settings. fis3 release
//fis.set('project.files', ['ExecFile']);
fis.set('project.ignore', ['ExecFile/**']);
// Global start


//启用插件 
fis.hook('relative'); 
//让所有文件，都使用相对路径。 
fis.match('**', { relative: true });

fis.match('app.less', {
	useHash: false,
  parser: fis.plugin('less'),
  rExt: '.css'
});


// default media is `dev`
fis.media('dev')
  .match('*', {
    useHash: false,
    optimizer: null
  });

// extends GLOBAL config
fis.media('production');