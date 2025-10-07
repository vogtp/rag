import { Component, Input, OnDestroy } from '@angular/core';
import { Subscription, filter, fromEvent } from 'rxjs';
import { ChatbotIcons } from './interfaces/library.interface';
import { ChatbotViewComponent } from './chatbot/chatbot.component';

@Component({
  selector: 'app-chatbot',
  standalone: true,
  imports: [ChatbotViewComponent],
  templateUrl: './chatbot.component.html',
  styleUrl: './chatbot.component.css',
})
export class ChatbotComponent implements OnDestroy {
  @Input({ required: true }) basePath!: string;
  @Input({ required: true }) model!: string;
  icons: ChatbotIcons = {
    chatbotIcon: '/static/icons/chatbot.svg',
    userIcon: '/static/icons/user.svg',
  };
  showTextBox: boolean = true;

  readonly keyDownEvent$ = fromEvent<KeyboardEvent>(document, 'keydown');
  private keyInputSub!: Subscription;

  ngOnDestroy(): void {
    if (this.keyInputSub) {
      this.keyInputSub.unsubscribe();
    }
  }

  //Method called whenever the chatbot icon is clicked
  onChatbotClicked(): void {
    this.showTextBox = true;
    this._subscribeToKeydownEvent();
  }

  onCloseChatbot(): void {
    this.showTextBox = false;
  }

  private _subscribeToKeydownEvent(): void {
    this.keyInputSub = this.keyDownEvent$
      .pipe(filter((event) => event.key === 'Escape'))
      .subscribe(() => {
        this.showTextBox = false;
        this.keyInputSub.unsubscribe();
      });
  }
}
