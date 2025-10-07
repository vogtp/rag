import { NgClass, NgStyle } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  ViewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TypingDirective } from '../directives/typing.directive';
import { ChatbotIcons } from '../interfaces/library.interface';
import {
  ChatCompletionMessage,
  ChatCompletionRequest,
  ChatCompletionResponse,
} from '../interfaces/openai.structs';

@Component({
  selector: 'lib-chatbot-textbox',
  standalone: true,
  imports: [NgStyle, NgClass, TypingDirective, FormsModule],
  templateUrl: './chatbot-textbox.component.html',
  styleUrl: './chatbot-textbox.component.css',
})
export class ChatbotTextboxComponent {
  @ViewChild('bodyChatbotContainer', { static: false })
  bodyContainer!: ElementRef;

  @Input({ required: true }) icons!: ChatbotIcons;
  @Input({ required: true }) basePath!: string;
  @Input({ required: true }) model!: string;
  @Output() closeChatbot = new EventEmitter<void>();

  readonly welcomeMessage: string = 'How can I help you?';
  readonly errorMessage: string = 'Something went wront. Please try later';
  inputText: string | undefined;
  waitingResponse: boolean = false;
  errorResponse: boolean = false;
  listOfMessages: ChatCompletionMessage[] = [
    //   { role: 'assistant', content: this.welcomeMessage },
  ];

  constructor(private http: HttpClient) {}

  getIcon(ind: number): string {
    return ind == 0 || ind % 2 == 0
      ? this.icons.chatbotIcon
      : this.icons.userIcon;
  }

  onCloseChatbot(): void {
    this.closeChatbot.emit();
  }

  onSendForm(): void {
    if (this.inputText != undefined) {
      let msg: ChatCompletionMessage = new ChatCompletionMessage();
      msg.role = 'user';
      msg.content = <string>this.inputText;
      this.listOfMessages.push(msg);
      this.inputText = undefined;
      this.waitingResponse = true;
      //Remove the welcome message
      let request: ChatCompletionRequest = new ChatCompletionRequest();
      request.model = this.model;
      request.messages = this.listOfMessages.slice(1);

      //Do the call
      this.http.post<ChatCompletionResponse>(this.basePath, request).subscribe({
        next: (res: ChatCompletionResponse) => {
          this.waitingResponse = false;
          if (res.choices) {
            let msg: ChatCompletionMessage = new ChatCompletionMessage();
            msg.role = res.choices[0].message.role;
            msg.content = res.choices[0].message.content;
            this.listOfMessages.push(msg);
          }
        },
        error: (err: any) => {
          this.waitingResponse = false;
          this.errorResponse = true;
          let msg: ChatCompletionMessage = new ChatCompletionMessage();
          msg.role = 'assistant';
          msg.content = <string>this.errorMessage + ': ' + err;
          this.listOfMessages.push();
        },
      });
    }
  }
}
