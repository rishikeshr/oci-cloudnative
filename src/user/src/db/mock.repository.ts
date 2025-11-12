import { Provider, Logger } from '@nestjs/common';
import { getRepositoryToken, getEntityManagerToken } from '@nestjs/typeorm';
import { Repository, EntityRepository, ObjectLiteral, EntityManager, getCustomRepository, ObjectType } from 'typeorm';

/**
 * Creates a mock repository for a TypeORM entity
 */
export class MockRepository extends Repository<any> {

  private static logger = new Logger(MockRepository.name);

  public static forEntity<T extends ObjectType<any>>(entity: T): Provider {
    return {
      provide: getRepositoryToken(entity as any),
      inject: [getEntityManagerToken()],
      useFactory: async (manager: EntityManager) => this.factory(manager, entity),
    };
  }

  private static factory<T extends ObjectType<any>>(manager: EntityManager, entity: T) {
    this.logger.warn(`Create mock repository for '${entity.name}'`);
    // tslint:disable-next-line:max-classes-per-file
    @EntityRepository(entity as any)
    class EntityRepo extends MockRepository { }
    return getCustomRepository(EntityRepo);
  }

  /**
   * Note: Removed manager property override for TypeScript 4.x compatibility.
   * TypeScript 4.x does not allow overriding a regular property with an accessor.
   * The base Repository.manager property is used instead.
   */

}
