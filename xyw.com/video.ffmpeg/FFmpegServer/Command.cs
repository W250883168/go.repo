using System;
using System.Diagnostics;
using System.Net;
using System.Text;
using System.Management;
namespace FFmpegServer
{
    /// <summary> 
    /// Command 的摘要说明。 
    /// </summary> 
    public class Command
    {
        //public Process proc = null;
        /// <summary> 
        /// 构造方法 
        /// </summary> 
        public Command()
        {
            //proc = new Process();
        }
        /// <summary> 
        /// 执行CMD语句,并且等待执行完成后返回
        /// </summary> 
        /// <param name="cmd">要执行的CMD命令</param> 
        public void RunCmd(string cmd)
        {
            Process proc = new Process();
            proc.StartInfo.CreateNoWindow = true;
            proc.StartInfo.FileName = "cmd.exe";
            proc.StartInfo.UseShellExecute = false;
            proc.StartInfo.RedirectStandardError = true;
            proc.StartInfo.RedirectStandardInput = true;
            proc.StartInfo.RedirectStandardOutput = true;
            proc.Start();
            proc.StandardInput.WriteLine(cmd);
            proc.WaitForExit();
            proc.Close();
        }
        /// <summary> 
        /// 打开软件并执行命令 
        /// </summary> 
        /// <param name="programName">软件路径加名称（.exe文件）</param> 
        /// <param name="cmd">要执行的命令</param> 
        public Process RunProgram(string programName, string cmd)
        {
            Process proc = new Process();
            proc.StartInfo.CreateNoWindow = true;
            proc.StartInfo.FileName = programName;
            proc.StartInfo.UseShellExecute = false;
            proc.StartInfo.RedirectStandardError = true;
            proc.StartInfo.RedirectStandardInput = true;
            proc.StartInfo.RedirectStandardOutput = true;
            proc.Start();
            if (cmd.Length != 0)
            {
                proc.StandardInput.WriteLine(cmd);
            }
            return proc;
        }

        /// <summary>
        /// 运行cmd命令，不显示命令窗口
        /// </summary>
        /// <param name="cmdExe">应用程序完整的路径</param>
        /// <param name="cmdStr">执行命令的参数</param>
        /// <returns></returns>
        public Process RunProgram2(string cmdExe, string cmddir, string cmdStr)
        {
            Process proc = new Process();
            try
            {
                ProcessStartInfo info = new ProcessStartInfo();
                info.FileName = cmdExe;
                info.RedirectStandardInput = true;
                info.RedirectStandardOutput = true;
                info.RedirectStandardError = true;
                info.Arguments = cmdStr;
                info.WindowStyle = ProcessWindowStyle.Hidden;
                info.CreateNoWindow = true;
                info.UseShellExecute = false;
                proc.StartInfo = info;
                proc.Start();
            }
            catch
            { }

            return proc;
        }

        public static void WriteLog(string context)
        {
            using (System.IO.StreamWriter sw = new System.IO.StreamWriter("C:\\statrlog.txt", true))
            {
                sw.WriteLine(context);
            }
        }
    }
}
