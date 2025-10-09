import { NgClass, NgStyle } from '@angular/common';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  ViewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { SbbLoadingIndicatorModule } from '@sbb-esta/angular/loading-indicator';
import { RemarkModule } from 'ngx-remark';
import {
  ChatCompletionMessage,
  ChatCompletionRequest,
  ChatCompletionResponse,
} from '../interfaces/openai.structs';

@Component({
  selector: 'chatbot',
  standalone: true,
  imports: [
    NgStyle,
    NgClass,
    FormsModule,
    MatInputModule,
    SbbLoadingIndicatorModule,
    RemarkModule,
  ],
  templateUrl: './chatbot.component.html',
  styleUrl: './chatbot.component.css',
})
export class ChatbotViewComponent {
  @ViewChild('bodyChatbotContainer', { static: false })
  bodyContainer!: ElementRef;

  @Input({ required: true }) basePath!: string;
  @Input({ required: true }) model!: string;
  @Output() closeChatbot = new EventEmitter<void>();
  inputText: string | undefined;
  waitingResponse: boolean = false;
  errorResponse: boolean = false;
  listOfMessages: ChatCompletionMessage[] = [  ];

  constructor(private http: HttpClient) {}

  onCloseChatbot(): void {
    this.closeChatbot.emit();
  }

  onSendForm(): void {
    console.log('Chat input: ' + this.inputText);
    if (this.inputText != undefined) {
      this.waitingResponse = true;
      let msg: ChatCompletionMessage = new ChatCompletionMessage();
      msg.role = 'user';
      msg.content = <string>this.inputText;
      this.listOfMessages.push(msg);
      this.inputText = undefined;
      //Remove the welcome message
      let request: ChatCompletionRequest = new ChatCompletionRequest();
      request.model = this.model;
      request.messages = this.listOfMessages; //.slice(1);

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
          console.log('Chat response error:' + err);
          console.log(err);
          this.waitingResponse = false;
          this.errorResponse = true;
          let msg: ChatCompletionMessage = new ChatCompletionMessage();
          msg.role = 'assistant';
          if (err instanceof HttpErrorResponse) {
            msg.content = err.error + ' (' + err.statusText + ')';
          } else {
            msg.content = err;
          }
          this.listOfMessages.push(msg);
        },
      });
    }
  }
}
