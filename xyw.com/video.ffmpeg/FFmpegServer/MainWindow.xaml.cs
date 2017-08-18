using System;
using System.Text;
using System.Windows;
using System.Windows.Forms;
using FFmpegService;


namespace FFmpegServer
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public UdpClientClass udpclient = null;
        private NotifyIcon notifyIcon;
        public string NumNo = "0";
        public MainWindow()
        {
            InitializeComponent();

            this.notifyIcon = new NotifyIcon();
            this.notifyIcon.BalloonTipText = "ffmpeg服务端运行中...";
            this.notifyIcon.ShowBalloonTip(2000);
            this.notifyIcon.Text = "ffmpeg服务端运行中...";
            //this.notifyIcon.Icon = new System.Drawing.Icon(@"favicon.ico");
            this.notifyIcon.Icon = System.Drawing.Icon.ExtractAssociatedIcon(System.Windows.Forms.Application.ExecutablePath);
            this.notifyIcon.Visible = true;
            //退出菜单项
            System.Windows.Forms.MenuItem exit = new System.Windows.Forms.MenuItem("退出");
            exit.Click += new EventHandler(Close);
            exit.Name = "exit";
            //关联托盘控件
            System.Windows.Forms.MenuItem[] childen = new System.Windows.Forms.MenuItem[] { exit };//open, minfrom, 
            notifyIcon.ContextMenu = new System.Windows.Forms.ContextMenu(childen);

            this.WindowState = System.Windows.WindowState.Maximized;
            this.Topmost = false;
            this.ShowInTaskbar = false;
            this.Visibility = System.Windows.Visibility.Hidden;
            this.Hide();

            string server = Properties.Settings.Default.ServerHost;
            int serverPort = Properties.Settings.Default.ServerPort;
            int thisPort = Properties.Settings.Default.ThisPort;
            udpclient = new UdpClientClass(server, serverPort, thisPort, 3000);

            CommandSendlog cmdlog = new CommandSendlog()
            {
                CmdIp = Properties.Settings.Default.ThisHost,
                CmdMac = Properties.Settings.Default.ThisMacAddr,
                CmdStr = "request",
                CmdPort = Properties.Settings.Default.ThisPort,
                CmdDate = DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss"),
                CmdType = "0",          // FFmpegServer
            };

            CmdMessage msg = new CmdMessage()
            {
                CmdID = DateTime.Now.ToString("yyyyMMdd-HHmmss"),
                CmdType = "CommandSendlog",
                JsonText = JsonHelper.Serialize(cmdlog)
            };

            string str = JsonHelper.Serialize(msg);
            byte[] senddata = Encoding.UTF8.GetBytes(str);
            udpclient.Send(senddata, senddata.Length);
        }

        private void Close(object sender, EventArgs e)
        {
            udpclient.Dispose();
            System.Windows.Application.Current.Shutdown();
        }
    }
}
