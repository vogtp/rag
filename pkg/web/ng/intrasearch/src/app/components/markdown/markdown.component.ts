import { Component, Input } from '@angular/core';
import { Marked } from 'marked';

@Component({
  selector: 'app-markdown',
  templateUrl: './markdown.component.html',
  styleUrls: ['./markdown.component.css'],
})
export class MarkdownComponent {
  processed: any = '';

  marked = new Marked();

  @Input() set markdown(markdown: string) {
    this.processed = this.marked.parse(markdown);
  }
}
