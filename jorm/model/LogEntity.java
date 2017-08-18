package com.hnlens.fingerchat.sso.dborm.model;
import java.util.Date;
public class LogEntity {
    private Long logID = 0L;
    private Integer logType = 0;
    private Date logDate;
    private String message = "";
    public Long getLogID() {
        return logID;
    }
    public void setLogID(Long logID) {
        this.logID = logID;
    }
    public Integer getLogType() {
        return logType;
    }
    public void setLogType(Integer logType) {
        this.logType = logType;
    }
    public Date getLogDate() {
        return logDate;
    }
    public void setLogDate(Date logDate) {
        this.logDate = logDate;
    }
    public String getMessage() {
        return message;
    }
    public void setMessage(String message) {
        this.message = message == null ? null : message.trim();
    }
}
