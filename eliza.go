// eliza.py - based on eliza.py. ported to Go by Tal Franji
// March, 2015. Original comment follows:
//----------------------------------------------------------------------
//  eliza.py
//
//  a cheezy little Eliza knock-off by Joe Strout <joe@strout.net>
//  with some updates by Jeff Epler <jepler@inetnebr.com>
//  hacked into a module and updated by Jez Higgins <jez@jezuk.co.uk>
//  last revised: 28 February 2005
//----------------------------------------------------------------------


package main
// msuserdog - Messaging Server - user entity - special user "Dog"
// Dog user mentors the user in using the app.
import (
  "bufio"
  "bytes"
  "fmt"
  "math/rand"
  "os"
  "regexp"
  "strconv"
  "strings"
)



var gReflections map[string]string = map[string]string{
  "am"   : "are",
  "was"  : "were",
  "i"    : "you",
  "i'd"  : "you would",
  "i've"  : "you have",
  "i'll"  : "you will",
  "my"  : "your",
  "are"  : "am",
  "you've": "I have",
  "you'll": "I will",
  "your"  : "my",
  "yours"  : "mine",
  "you"  : "me",
  "me"  : "you",
}

func reflectWords(s string) string {
  words := regexp.MustCompile(`\s+`).Split(s, -1)
  var buf bytes.Buffer
  for i, word := range(words) {
    if i > 0 {
      buf.WriteString(" ")
    }
    newWord, ok := gReflections[strings.ToLower(word)]
    if ok {
      buf.WriteString(newWord)
    } else {
      buf.WriteString(word)
    }
  }
  return buf.String()
}

type ElizaRule struct {
  Pat string
  Values []string
}

var gRules []ElizaRule = []ElizaRule{
  {"I need (.*)",
  []string{  "Why do you need %1?",
    "Would it really help you to get %1?",
    "Are you sure you need %1?"}},
  
  {`Why don\'?t you ([^\?]*)\??`,
  []string{  "Do you really think I don't %1?",
    "Perhaps eventually I will %1.",
    "Do you really want me to %1?"}},
  
  {`Why can\'?t I ([^\?]*)\??`,
  []string{  "Do you think you should be able to %1?",
    "If you could %1, what would you do?",
    "I don't know -- why can't you %1?",
    "Have you really tried?"}},
    
  {`Who (.*)`,
  []string{ "Three gueses...",
    "Your best friend.",
    "It's nobody."}},
    
  {".*[\u05d0-\u05ea]+",
  []string{ "Sorry, I don't speak Hebrew.",
    "Can you write in English, please?",
    "Is that Greek?"}},
  
  {`I can\'?t (.*)`,
  []string{  "How do you know you can't %1?",
    "Perhaps you could %1 if you tried.",
    "What would it take for you to %1?"}},
  
  {"I am (.*)",
  []string{  "Did you come to me because you are %1?",
    "How long have you been %1?",
    "How do you feel about being %1?"}},
  
  {`I\'?m (.*)`,
  []string{  "How does being %1 make you feel?",
    "Do you enjoy being %1?",
    "Why do you tell me you're %1?",
    "Why do you think you're %1?"}},
  
  {`Are you ([^\?]*)\??`,
  []string{  "Why does it matter whether I am %1?",
    "Would you prefer it if I were not %1?",
    "Perhaps you believe I am %1.",
    "I may be %1 -- what do you think?"}},
  
  {"What (.*)",
  []string{  "Why do you ask?",
    "How would an answer to that help you?",
    "What do you think?"}},
  
  {"How (.*)",
  []string{  "How do you suppose?",
    "Perhaps you can answer your own question.",
    "What is it you're really asking?"}},
  
  {"Because (.*)",
  []string{  "Is that the real reason?",
    "What other reasons come to mind?",
    "Does that reason apply to anything else?",
    "If %1, what else must be true?"}},
  
  {"(.*) sorry (.*)",
  []string{  "There are many times when no apology is needed.",
    "What feelings do you have when you apologize?"}},
  
  {"(?:Hello|Hi|Hey)(.*)",
  []string{  "Hello... I'm glad you could drop by today.",
    "Hi there... how are you today?",
    "Hello, how are you feeling today?"}},
  
  {"I think (.*)",
  []string{  "Do you doubt %1?",
    "Do you really think so?",
    "But you're not sure %1?"}},
  
  {"(.*) friend (.*)",
  []string{  "Tell me more about your friends.",
    "When you think of a friend, what comes to mind?",
    "Why don't you tell me about a childhood friend?"}},
  
  {"Yes",
  []string{  "You seem quite sure.",
    "OK, but can you elaborate a bit?"}},

  {"(.*) computer(.*)",
  []string{  "Are you really talking about me?",
    "Does it seem strange to talk to a computer?",
    "How do computers make you feel?",
    "Do you feel threatened by computers?"}},

  {"Is it (.*)",
  []string{  "Do you think it is %1?",
    "Perhaps it's %1 -- what do you think?",
    "If it were %1, what would you do?",
    "It could well be that %1."}},
  
  {`It is (.*)`,
  []string{  "You seem very certain.",
    "If I told you that it probably isn't %1, what would you feel?"}},
  
  {`Can you ([^\?]*)\??`,
  []string{  "What makes you think I can't %1?",
    "If I could %1, then what?",
    "Why do you ask if I can %1?"}},
  
  {`Can I ([^\?]*)\??`,
  []string{  "Perhaps you don't want to %1.",
    "Do you want to be able to %1?",
    "If you could %1, would you?"}},
  
  {"You are (.*)",
  []string{  "Why do you think I am %1?",
    "Does it please you to think that I'm %1?",
    "Perhaps you would like me to be %1.",
    "Perhaps you're really talking about yourself?"}},
  
  {`You\'?re (.*)`,
  []string{  "Why do you say I am %1?",
    "Why do you think I am %1?",
    "Are we talking about you, or me?"}},
  
  {`I don\'?t (.*)`,
  []string{  "Don't you really %1?",
    "Why don't you %1?",
    "Do you want to %1?"}},
  
  {"I feel (.*)",
  []string{  "Good, tell me more about these feelings.",
    "Do you often feel %1?",
    "When do you usually feel %1?",
    "When you feel %1, what do you do?"}},
  
  {"I have (.*)",
  []string{  "Why do you tell me that you've %1?",
    "Have you really %1?",
    "Now that you have %1, what will you do next?"}},
  
  {"I would (.*)",
  []string{  "Could you explain why you would %1?",
    "Why would you %1?",
    "Who else knows that you would %1?"}},
  
  {"Is there (.*)",
  []string{  "Do you think there is %1?",
    "It's likely that there is %1.",
    "Would you like there to be %1?"}},
  
  {"My (.*)",
  []string{  "I see, your %1.",
    "Why do you say that your %1?",
    "When your %1, how do you feel?"}},
  
  {"You (.*)",
  []string{  "We should be discussing you, not me.",
    "Why do you say that about me?",
    "Why do you care whether I %1?"}},
    
  {"Why (.*)",
  []string{  "Why don't you tell me the reason why %1?",
    "Why do you think %1?" }},
    
  {"I want (.*)",
  []string{  "What would it mean to you if you got %1?",
    "Why do you want %1?",
    "What would you do if you got %1?",
    "If you got %1, then what would you do?"}},
  
  {"(.*) mother(.*)",
  []string{  "Tell me more about your mother.",
    "What was your relationship with your mother like?",
    "How do you feel about your mother?",
    "How does this relate to your feelings today?",
    "Good family relations are important."}},
  
  {"(.*) father(.*)",
  []string{  "Tell me more about your father.",
    "How did your father make you feel?",
    "How do you feel about your father?",
    "Does your relationship with your father relate to your feelings today?",
    "Do you have trouble showing affection with your family?"}},

  {"(.*) child(.*)",
  []string{  "Did you have close friends as a child?",
    "What is your favorite childhood memory?",
    "Do you remember any dreams or nightmares from childhood?",
    "Did the other children sometimes tease you?",
    "How do you think your childhood experiences relate to your feelings today?"}},
    
  {`(.*)\?`,
  []string{  "Why do you ask that?",
    "Please consider whether you can answer your own question.",
    "Perhaps the answer lies within yourself?",
    "Why don't you tell me?"}},
  
  {"quit",
  []string{  "Thank you for talking with me.",
    "Good-bye.",
    "Thank you, that will be $150.  Have a good day!"}},
  
  {"(.*)",
  []string{  "Please tell me more.",
    "Let's change focus a bit... Tell me about your family.",
    "Can you elaborate on that?",
    "Why do you say that %1?",
    "I see.",
    "Very interesting.",
    "%1.",
    "I see.  And what does that tell you?",
    "How does that make you feel?",
    "How do you feel when you say that?"}},
  }

func ElizaRespond(s string) string {
  for _, rule := range(gRules) {
    // m is the array of sub-matches
    m := regexp.MustCompile(strings.ToLower(rule.Pat)).FindStringSubmatch(s)
    if len(m) == 0 {
      continue
    }
    // found a match ... stuff with corresponding value
    // chosen randomly from among the available options
    i_resp := rand.Intn(len(rule.Values))
    resp := rule.Values[i_resp]
    // we've got a response... stuff in reflected text where indicated
    pos := strings.Index(resp,"%")
    for pos > -1 {
      num, _ := strconv.ParseInt(resp[pos+1:pos+2], 10, 8)
      resp = resp[:pos] + reflectWords(m[num]) + resp[pos+2:]
      pos = strings.Index(resp,"%")
    }
    // fix munged punctuation at the end
    if resp[len(resp) - 2:] == "?." {
       resp = resp[:len(resp) - 2] + "."
    }
    if resp[len(resp) - 2:] == "??" {
       resp = resp[:len(resp) - 2] + "?"
    }
    return resp
  }
  return "Ha?"
}


func main() {
  reader := bufio.NewReader(os.Stdin)
  fmt.Printf("Therapist\n---------\n")
  fmt.Printf("Talk to the program by typing in plain English, using normal upper-\n")
  fmt.Printf("and lower-case letters and punctuation.  Enter \"quit\" when done.\n")
  fmt.Printf("=======================================================================================\n")
  fmt.Printf("Hello.  How are you feeling today?\n")
  s := ""
  for  {
    if s == "quit" {
      break
    }
    fmt.Printf(">")
    s, _ = reader.ReadString('\n')
    s = strings.ToLower(strings.TrimRight(s, "\n\r."))
    fmt.Printf("%s\n", ElizaRespond(s))
  }

}