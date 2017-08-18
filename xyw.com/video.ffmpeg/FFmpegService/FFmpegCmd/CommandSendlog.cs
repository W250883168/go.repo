using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace FFmpegService
{
    public class CommandSendlog
    {
        public int Id { get; set; }
        public int Classroomid { get; set; }
        public string CmdIp { get; set; }
        public int CmdPort { get; set; }
        public string CmdStr { get; set; }
        public string CmdType { get; set; }
        public int CmdUsersId { get; set; }
        public string CmdUsersName { get; set; }
        public string CmdDate { get; set; }
        public int CmdState { get; set; }
        public string CmdError { get; set; }
        public string CmdMac { get; set; }

        public CommandSendlog()
        {
            this.CmdIp = "";
            this.CmdStr = "";
            this.CmdType = "";
            this.CmdUsersName = "";
            this.CmdDate = "";
            this.CmdError = "";
            this.CmdMac = "";
        }
    }
}
