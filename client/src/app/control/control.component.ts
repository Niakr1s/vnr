import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { SentenceService } from '../services/sentence.service';

@Component({
  selector: 'app-control',
  templateUrl: './control.component.html',
  styleUrls: ['./control.component.css'],
})
export class ControlComponent implements OnInit, OnDestroy {
  total!: number;
  private subs: Subscription[] = [];

  constructor(private sentenceService: SentenceService) {}

  ngOnInit(): void {
    this.subs.push(
      this.sentenceService.totalSentences$.subscribe({
        next: (total) => (this.total = total),
      })
    );
  }

  ngOnDestroy(): void {
    this.subs.forEach((s) => s.unsubscribe());
  }

  onDelete(): void {
    this.sentenceService.deleteCurrentSentence();
  }

  onNext(): void {
    this.sentenceService.next();
  }

  onPrev(): void {
    this.sentenceService.prev();
  }

  onLast(): void {
    this.sentenceService.last();
  }

  hasNext(): boolean {
    return this.sentenceService.hasNext();
  }

  hasPrev(): boolean {
    return this.sentenceService.hasPrev();
  }
}
