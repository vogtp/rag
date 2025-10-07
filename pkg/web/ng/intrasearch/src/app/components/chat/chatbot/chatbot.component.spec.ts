import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChatbotViewComponent } from './chatbot.component';

describe('ChatbotTextboxComponent', () => {
  let component: ChatbotViewComponent;
  let fixture: ComponentFixture<ChatbotViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ChatbotViewComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ChatbotViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
