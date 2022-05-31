package main

import(
  "net/http"
  "io/ioutil"
  "encoding/json"
  "fmt"
  "log"
  "bytes"
  "strconv"
  "youtube"
)


func main(){
  botToken:="5015864716:AAGjVQyrEV0ORoKjmmC976y91DJiinv79OA"
  botApi:= "https://api.telegram.org/bot"
  botUrl:=botApi + botToken
  offset := 0
  
  for ;; {
    updates, err:=getUpdates(botUrl,offset)
    if err != nil{
      log.Println("ERR", err.Error())
    }
    
    for _, update := range updates{
      err=respond(botUrl,update)
      offset = update.UpdateId + 1
    }
    fmt.Println(updates)
  }
}

func getUpdates(botUrl string,offset int)([]Update, error){
  resp, err:= http.Get(botUrl+"/getUpdates" + "?offset=" + strconv.Itoa(offset))
  if err != nil{
    return nil, err
  }
  defer resp.Body.Close()
  
  body, err:=ioutil.ReadAll(resp.Body)
  if err != nil{
    return nil, err
  }
  
  var restResponse RestResponse
  err=json.Unmarshal(body, &restResponse)
  if err !=nil{
    return nil, err
  }
  return restResponse.Result, nil
}

func respond(botUrl string,update Update) (error){
  var botMessage BotMessage
  botMessage.ChatId = update.Message.Chat.ChatId
  videoUrl, err := youtube.GetLastVideo(update.Message.Text)
  
  botMessage.Text = videoUrl
  if err !=nil{
    return err
  }

   buf,err:=json.Marshal(botMessage)
  if err !=nil{
    return err
  }
  
  _, err=http.Post(botUrl + "/sendMessage", "application/json", bytes.NewBuffer(buf))
  if err !=nil{
    return err
  }
  return nil
}
