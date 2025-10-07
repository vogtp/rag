import { Component, Input } from '@angular/core';
import { MatTabsModule } from '@angular/material/tabs';
import { ChatbotComponent } from '../chat/chatbot.component';
import { SearchComponent } from '../search/search.component';

@Component({
  selector: 'app-collection-overview',
  standalone: true,
  imports: [ChatbotComponent, MatTabsModule, SearchComponent],
  templateUrl: './collection-overview.component.html',
  styleUrl: './collection-overview.component.css',
})
export class CollectionOverviewComponent {
  @Input() collection: string = 'intranet-all';

  basePath: string = 'http://localhost:4444/api/chat/completions';
}
