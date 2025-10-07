import { Component, Input } from '@angular/core';
import { MatTabsModule } from '@angular/material/tabs';
import { ChatbotComponent } from '../chat/chatbot.component';
import { ChatbotIcons } from '../chat/interfaces/library.interface';
import { SearchComponent } from '../search/search.component';

@Component({
  selector: 'app-collection',
  imports: [ChatbotComponent, MatTabsModule, SearchComponent],
  templateUrl: './collection.component.html',
  styleUrl: './collection.component.css',
})
export class CollectionComponent {
  @Input() collection: string = 'intranet-all';

  
  basePath: string = 'http://localhost:4444/api/chat/completions';
}
