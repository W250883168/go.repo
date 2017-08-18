package com.hnlens.fingerchat.sso.dborm.model;
public class ApplicationEntity {
    private String appID = "";
    private Integer appType = 0;
    private String appName = "";
    private String homeURL = "";
    private String logoutURL = "";
    private Integer appState = 0;
    public String getAppID() {
        return appID;
    }
    public void setAppID(String appID) {
        this.appID = appID == null ? null : appID.trim();
    }
    public Integer getAppType() {
        return appType;
    }
    public void setAppType(Integer appType) {
        this.appType = appType;
    }
    public String getAppName() {
        return appName;
    }
    public void setAppName(String appName) {
        this.appName = appName == null ? null : appName.trim();
    }
    public String getHomeURL() {
        return homeURL;
    }
    public void setHomeURL(String homeURL) {
        this.homeURL = homeURL == null ? null : homeURL.trim();
    }
    public String getLogoutURL() {
        return logoutURL;
    }
    public void setLogoutURL(String logoutURL) {
        this.logoutURL = logoutURL == null ? null : logoutURL.trim();
    }
    public Integer getAppState() {
        return appState;
    }
    public void setAppState(Integer appState) {
        this.appState = appState;
    }
}
